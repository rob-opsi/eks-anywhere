package factory_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/cluster"
	mockswriter "github.com/aws/eks-anywhere/pkg/filewriter/mocks"
	clientProviders "github.com/aws/eks-anywhere/pkg/providers/aws/mocks"
	"github.com/aws/eks-anywhere/pkg/providers/docker"
	dockerMocks "github.com/aws/eks-anywhere/pkg/providers/docker/mocks"
	"github.com/aws/eks-anywhere/pkg/providers/factory"
	"github.com/aws/eks-anywhere/pkg/providers/vsphere"
	vsphereMocks "github.com/aws/eks-anywhere/pkg/providers/vsphere/mocks"
	releasev1alpha1 "github.com/aws/eks-anywhere/release/api/v1alpha1"
)

// makes sure right type of provider is created based on the input
func TestProviderFactoryBuildProvider(t *testing.T) {
	type providerMatch struct {
		kind    string
		version string
	}
	type args struct {
		clusterConfigFileName string
		clusterConfig         *v1alpha1.Cluster
	}
	clusterSpec := &cluster.Spec{
		VersionsBundle: &cluster.VersionsBundle{
			VersionsBundle: &releasev1alpha1.VersionsBundle{
				VSphere: releasev1alpha1.VSphereBundle{Version: "v0.7.8"},
				Docker:  releasev1alpha1.DockerBundle{Version: "v0.3.19"},
			},
		},
	}
	tests := []struct {
		name    string
		args    args
		want    providerMatch
		wantErr error
	}{
		{
			name: "Vsphere cluster",
			args: args{
				clusterConfig: &v1alpha1.Cluster{Spec: v1alpha1.ClusterSpec{
					DatacenterRef: v1alpha1.Ref{
						Kind: v1alpha1.VSphereDatacenterKind,
					},
				}},
				clusterConfigFileName: "testdata/cluster_vsphere.yaml",
			},
			want: providerMatch{
				kind:    vsphere.ProviderName,
				version: "v0.7.8",
			},
		},
		{
			name: "Docker cluster",
			args: args{
				clusterConfig: &v1alpha1.Cluster{Spec: v1alpha1.ClusterSpec{
					DatacenterRef: v1alpha1.Ref{
						Kind: v1alpha1.DockerDatacenterKind,
					},
				}},
				clusterConfigFileName: "testdata/cluster_docker.yaml",
			},
			want: providerMatch{
				kind:    docker.ProviderName,
				version: "v0.3.19",
			},
		},
		{
			name: "Aws cluster not supported",
			args: args{
				clusterConfig: &v1alpha1.Cluster{Spec: v1alpha1.ClusterSpec{
					DatacenterRef: v1alpha1.Ref{
						Kind: v1alpha1.AWSDatacenterKind,
					},
				}},
				clusterConfigFileName: "testdata/cluster_aws.yaml",
			},
			wantErr: fmt.Errorf("valid providers include: %s, %s", docker.ProviderName, vsphere.ProviderName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			mockCtrl := gomock.NewController(t)
			p := &factory.ProviderFactory{
				AwsClient:            clientProviders.NewMockProviderClient(mockCtrl),
				DockerClient:         dockerMocks.NewMockProviderClient(mockCtrl),
				VSphereGovcClient:    vsphereMocks.NewMockProviderGovcClient(mockCtrl),
				VSphereKubectlClient: vsphereMocks.NewMockProviderKubectlClient(mockCtrl),
				Writer:               mockswriter.NewMockFileWriter(mockCtrl),
			}
			got, err := p.BuildProvider(tt.args.clusterConfigFileName, tt.args.clusterConfig)
			if err == nil {
				if got.Name() != tt.want.kind || got.Version(clusterSpec) != tt.want.version {
					t.Errorf("BuildProvider() got = %v, want %v", got, tt.want)
				}
			}

			if tt.wantErr != nil {
				g.Expect(err).To(Equal(tt.wantErr))
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
