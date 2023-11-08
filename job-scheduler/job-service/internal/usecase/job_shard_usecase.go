package usecase

// JobShardUsecase is an interface that represents the job shard usecase
type JobShardUsecase interface {
	GetSharedId() uint64
}
