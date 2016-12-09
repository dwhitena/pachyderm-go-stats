package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/pachyderm/pachyderm/src/client"
	"github.com/pachyderm/pachyderm/src/client/pps"
)

// projects includes a list of example Go projects on Github.
var projects = []string{
	"docker/docker",
	"kubernetes/kubernetes",
	"hashicorp/consul",
	"spf13/hugo",
	"prometheus/prometheus",
	"influxdata/influxdb",
	"coreos/etcd",
}

func main() {

	// Connect to Pachyderm.
	c, err := client.NewFromAddress("0.0.0.0:30650")
	if err != nil {
		log.Fatal(err)
	}

	// Create a repo called "projects."
	if err := c.CreateRepo("projects"); err != nil {
		log.Fatal(err)
	}

	// Loop over the projects.
	for idx, project := range projects {

		// Start a commit in our "projects" data repo on the "master" branch.
		commit, err := c.StartCommit("projects", "master")
		if err != nil {
			log.Fatal(err)
		}

		// Put a file containing the respective project name.
		if _, err := c.PutFile("projects", commit.ID, strconv.Itoa(idx), strings.NewReader(project)); err != nil {
			log.Fatal(err)
		}

		// Finish the commit.
		if err := c.FinishCommit("projects", commit.ID); err != nil {
			log.Fatal(err)
		}
	}

	// Define the stdin for our pipeline.
	stdin := []string{
		"for filename in /pfs/projects/*; do",
		"REPONAME=`cat $filename`",
		"source /stats.sh >> /pfs/out/results",
		"done",
	}

	// Create our stats pipeline.
	if err := c.CreatePipeline(
		"stats",               // the name of the pipeline
		"dwhitena/stats",      // your docker image
		[]string{"/bin/bash"}, // the command run in your docker image
		stdin,
		nil, // let pachyderm decide the parallelism
		[]*pps.PipelineInput{
			// reduce over "projects"
			client.NewPipelineInput("projects", client.ReduceMethod),
		},
		false, // not an update
	); err != nil {
		log.Fatal(err)
	}
}
