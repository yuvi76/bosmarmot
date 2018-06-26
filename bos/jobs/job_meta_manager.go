package jobs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monax/bosmarmot/bos/definitions"
	"github.com/monax/bosmarmot/bos/loader"
	log "github.com/sirupsen/logrus"
)

func MetaJob(meta *definitions.Meta, do *definitions.Packages) (string, error) {
	var err error
	var pwd string

	pwd, err = os.Getwd()
	if err != nil {
		return "failed", err
	}

	// work from a fresh Do object
	newDo := definitions.NewPackage()
	newDo.Address = do.Address
	newDo.ChainURL = do.ChainURL
	newDo.CurrentOutput = do.CurrentOutput
	newDo.DefaultAmount = do.DefaultAmount
	newDo.DefaultFee = do.DefaultFee
	newDo.DefaultGas = do.DefaultGas
	newDo.DefaultSets = do.DefaultSets
	newDo.PublicKey = do.PublicKey
	newDo.Signer = do.Signer

	// Set subYAMLPath
	newDo.YAMLPath = meta.File

	// if subYAMLPath does not exist, try YAMLPath relative to do.Path
	if _, err := os.Stat(newDo.YAMLPath); os.IsNotExist(err) {
		newDo.YAMLPath = filepath.Join(do.Path, newDo.YAMLPath)
	}

	// if subYAMLPath does not exist, try YAMLPath relative to pwd
	if _, err := os.Stat(newDo.YAMLPath); os.IsNotExist(err) {
		newDo.YAMLPath = filepath.Join(pwd, newDo.YAMLPath)
	}

	// if subYAMLPath cannot be found, abort
	if _, err := os.Stat(newDo.YAMLPath); os.IsNotExist(err) {
		return "failed", fmt.Errorf("could not find sub YAML file (%s)",
			do.YAMLPath)
	}

	// once we have the proper subYAMLPath set the paths accordingly
	newDo.Path = filepath.Dir(newDo.YAMLPath)
	newDo.BinPath = filepath.Join(newDo.Path, filepath.Base(do.BinPath))
	newDo.ABIPath = filepath.Join(newDo.Path, filepath.Base(do.ABIPath))

	// load the package
	log.WithField("=>", newDo.YAMLPath).Info("Loading sub YAML")
	newDo.Package, err = loader.LoadPackage(newDo.YAMLPath)
	if err != nil {
		return "failed", err
	}

	// set the deploy contract jobs relative to the newDo's root directory
	for _, job := range newDo.Package.Jobs {
		if job.Deploy != nil {
			job.Deploy.Contract = filepath.Join(newDo.Path, job.Deploy.Contract)
		}
	}

	err = RunJobs(newDo)
	if err != nil {
		return "failed", err
	}

	do.CurrentOutput = ""
	return "passed", nil
}