package testutils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kalensk/plusy/src/db/redisdb"
	"github.com/kalensk/plusy/src/messages"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var SekretChatRoom = messages.Chat{ID: 789, Title: "sekret-room"}
var UserDougCat = &messages.User{ID: 1, FirstName: "doug", Username: "doug_cat"}
var UserDougBat = &messages.User{ID: 2, FirstName: "doug", Username: "doug_bat"}
var UserMrRogers = &messages.User{ID: 3, FirstName: "mr", Username: "mr_rogers"}

func StringPointer(input string) *string {
	return &input
}

func Int64Pointer(input int64) *int64 {
	return &input
}

type Docker struct {
	m        *testing.M
	Log      *logrus.Logger
	Database *redisdb.Redis
	options  *dockertest.RunOptions
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func NewRedisDocker(m *testing.M) *Docker {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.Formatter.(*logrus.TextFormatter).DisableLevelTruncation = false

	options := &dockertest.RunOptions{Repository: "redis", Tag: "5.0-rc-stretch", Name: "test-redis"}

	dockerContainer := &Docker{m: m, Log: log, options: options}
	return dockerContainer
}

func New(m *testing.M, repository string, tag string, containerName string) *Docker {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.Formatter.(*logrus.TextFormatter).DisableLevelTruncation = false

	options := &dockertest.RunOptions{Repository: repository, Tag: tag, Name: containerName}

	return &Docker{m: m, Log: log, options: options}
}

func (d *Docker) removeRunningContainer(pool *dockertest.Pool) error {
	runningContainers, _ := pool.Client.ListContainers(docker.ListContainersOptions{All: true})

	for _, container := range runningContainers {
		for _, name := range container.Names {
			if strings.Contains(name, d.options.Name) {
				d.Log.Infof("removing previously running container: %s", name)
				if err := pool.Client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, RemoveVolumes: true, Force: true}); err != nil {
					d.Log.Fatal("failed to remove previously running container: %s", d.options.Name)
				}
				return nil
			}
		}
	}

	return errors.Errorf("no previously running containers found with name: %s", d.options.Name)
}

func (d *Docker) SetupTestRedis() *redisdb.Redis {
	d.Log.Info("test setup")
	d.Log.Info("creating test redis docker resource")
	pool, err := dockertest.NewPool("")
	d.pool = pool
	if err != nil {
		d.Log.Fatalf("could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(d.options)
	d.resource = resource
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "container already exists") {
			d.Log.Warnf("container %s already exists: attempting to remove and recreate", d.options.Name)
			err = d.removeRunningContainer(pool)
			if err != nil {
				d.Log.Fatal("could not remove previously running container")
			}
			return d.SetupTestRedis()
		}

		d.Log.Fatal("failed to start redis test resource")
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	retryCallback := func() error {
		d.Database, err = redisdb.New(d.Log, fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")))
		if err != nil {
			return err
		}
		return d.Database.Ping()
	}

	if err := pool.Retry(retryCallback); err != nil {
		d.Log.Fatalf("failed to connect to docker container: %s", err)
	}

	return d.Database
}

func (d *Docker) GetDatabaseConnectionString() string {
	return "localhost:" + d.resource.GetPort("6379/tcp")
}

func (d *Docker) Run() int {
	d.Log.Info("running tests")
	returnCode := d.m.Run()
	return returnCode
}

func (d *Docker) TearDown(returnCode int) {
	d.Log.Info("test teardown")
	d.Log.Info("cleaning up test redis resources")
	err := d.pool.Purge(d.resource)
	if err != nil {
		d.Log.Fatal("could not purge resource: %s", err)
	}

	os.Exit(returnCode)
}
