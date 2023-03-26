## Description

This is an internal tool to improve the time it takes to adapt a new jacob batch
to run on a kubernetes cluster.

This cli wil generate a kustomize hierarchy for you, and will generate two environments folder,
one for coll (quality) and one for prod (production), and will generate a kustomization.yaml.
Then the program copies the property files and edit the datasource property and the base path folder property
to point to the correct environment.

## Usage

Install the program with
```bash
go install
```

Then into the folder where you want to generate the kustomize hierarchy, run 
```bash
jacobCloudAdapert adapt --programName="name"
```

