package module

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindModules(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		repoPath        string
		expectedModules Modules
		expectedErr     error
	}{
		{
			"testdata/find_modules/terraform-aws-eks",
			Modules{
				&Module{
					title:       "ALB Ingress Controller Module",
					description: "This Terraform Module installs and configures the [AWS ALB Ingress Controller](https://github.com/kubernetes-sigs/aws-alb-ingress-controller) on an EKS cluster.",
					url:         "https:/github.com/gruntwork-io/terraform-aws-eks/tree/master/modules/eks-alb-ingress-controller",
					path:        "modules/eks-alb-ingress-controller",
				},
				&Module{
					title:       "ALB Ingress Controller IAM Policy Module",
					description: "This Terraform Module defines an [IAM policy](http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/QuickStartEC2Instance.html#d0e22325)  that defines the minimal set of permissions necessary for the [AWS ALB Ingress Controller]",
					url:         "https:/github.com/gruntwork-io/terraform-aws-eks/tree/master/modules/eks-alb-ingress-controller-iam-policy",
					path:        "modules/eks-alb-ingress-controller-iam-policy",
				},
				&Module{
					title:       "EKS AWS Auth Merger",
					description: "This module contains a go CLI, docker container, and terraform module for deploying a Kubernetes controller for managing mappings between AWS IAM roles and users to RBAC groups in Kubernetes. The official way to manage the mapping is to add values in a single, central `ConfigMap`.  This module allows you to break up the central `ConfigMap` across multiple.   toc::[]",
					url:         "https:/github.com/gruntwork-io/terraform-aws-eks/tree/master/modules/eks-aws-auth-merger",
					path:        "modules/eks-aws-auth-merger",
				}},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.repoPath, func(t *testing.T) {
			// Unfortunately, we are unable to commit the `.git` directory. We have to temporarily rename it while running the tests.
			os.Rename(filepath.Join(testCase.repoPath, "gitdir"), filepath.Join(testCase.repoPath, ".git"))
			defer os.Rename(filepath.Join(testCase.repoPath, ".git"), filepath.Join(testCase.repoPath, "gitdir"))

			ctx := context.Background()
			modules, err := FindModules(ctx, testCase.repoPath)

			for _, module := range modules {
				currentDir, err := os.Getwd()
				assert.NoError(t, err)

				relPath, err := filepath.Rel(filepath.Join(currentDir, testCase.repoPath), module.path)
				assert.NoError(t, err)

				module.path = relPath
				module.content = ""
			}

			assert.Equal(t, testCase.expectedModules, modules)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}

}
