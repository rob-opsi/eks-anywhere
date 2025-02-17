package interfaces

import (
	"context"

	"github.com/aws/eks-anywhere/pkg/bootstrapper"
	"github.com/aws/eks-anywhere/pkg/cluster"
	"github.com/aws/eks-anywhere/pkg/providers"
	"github.com/aws/eks-anywhere/pkg/types"
	"github.com/aws/eks-anywhere/pkg/validations"
)

type Bootstrapper interface {
	CreateBootstrapCluster(ctx context.Context, clusterSpec *cluster.Spec, opts ...bootstrapper.BootstrapClusterOption) (*types.Cluster, error)
	DeleteBootstrapCluster(context.Context, *types.Cluster, bool) error
}

type ClusterManager interface {
	MoveCAPI(ctx context.Context, from, to *types.Cluster, checkers ...types.NodeReadyChecker) error
	CreateWorkloadCluster(ctx context.Context, managementCluster *types.Cluster, clusterSpec *cluster.Spec, provider providers.Provider) (*types.Cluster, error)
	UpgradeCluster(ctx context.Context, managementCluster, workloadCluster *types.Cluster, clusterSpec *cluster.Spec, provider providers.Provider) error
	DeleteCluster(ctx context.Context, managementCluster, clusterToDelete *types.Cluster) error
	InstallCAPI(ctx context.Context, clusterSpec *cluster.Spec, cluster *types.Cluster, provider providers.Provider) error
	InstallNetworking(ctx context.Context, cluster *types.Cluster, clusterSpec *cluster.Spec) error
	InstallStorageClass(ctx context.Context, cluster *types.Cluster, provider providers.Provider) error
	SaveLogs(ctx context.Context, cluster *types.Cluster) error
	InstallCustomComponents(ctx context.Context, clusterSpec *cluster.Spec, cluster *types.Cluster) error
	CreateEKSAResources(ctx context.Context, cluster *types.Cluster, clusterSpec *cluster.Spec, datacenterConfig providers.DatacenterConfig, machineConfigs []providers.MachineConfig) error
	PauseEKSAControllerReconcile(ctx context.Context, cluster *types.Cluster, clusterSpec *cluster.Spec, provider providers.Provider) error
	ResumeEKSAControllerReconcile(ctx context.Context, cluster *types.Cluster, clusterSpec *cluster.Spec, provider providers.Provider) error
	EKSAClusterSpecChanged(ctx context.Context, cluster *types.Cluster, clusterSpec *cluster.Spec, datacenterConfig providers.DatacenterConfig, machineConfigs []providers.MachineConfig) (bool, error)
	InstallMachineHealthChecks(ctx context.Context, workloadCluster *types.Cluster, provider providers.Provider) error
}

type AddonManager interface {
	InstallGitOps(ctx context.Context, cluster *types.Cluster, clusterSpec *cluster.Spec, datacenterConfig providers.DatacenterConfig, machineConfigs []providers.MachineConfig) error
	PauseGitOpsKustomization(ctx context.Context, cluster *types.Cluster, clusterSpec *cluster.Spec) error
	ResumeGitOpsKustomization(ctx context.Context, cluster *types.Cluster, clusterSpec *cluster.Spec) error
	UpdateGitEksaSpec(ctx context.Context, clusterSpec *cluster.Spec, datacenterConfig providers.DatacenterConfig, machineConfigs []providers.MachineConfig) error
	ForceReconcileGitRepo(ctx context.Context, cluster *types.Cluster, clusterSpec *cluster.Spec) error
	Validations(ctx context.Context, clusterSpec *cluster.Spec) []validations.Validation
	CleanupGitRepo(ctx context.Context, clusterSpec *cluster.Spec) error
}

type Validator interface {
	PreflightValidations(ctx context.Context) error
}
