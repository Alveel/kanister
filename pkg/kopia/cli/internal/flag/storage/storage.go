// Copyright 2024 The Kanister Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"sync"

	cmdlog "github.com/kanisterio/kanister/pkg/kopia/cli/internal/log"
	"github.com/kanisterio/kanister/pkg/log"
	rs "github.com/kanisterio/kanister/pkg/secrets/repositoryserver"

	"github.com/kanisterio/kanister/pkg/kopia/cli/internal/flag/storage/fs"
	"github.com/kanisterio/kanister/pkg/kopia/cli/internal/flag/storage/model"
)

// Option is a function that sets a storage option.
type Option func(*model.StorageFlag)

// WithLogger sets the logger for the storage.
func WithLogger(logger log.Logger) Option {
	return func(s *model.StorageFlag) {
		s.Logger = logger
	}
}

// WithFactory sets the storage args builder factory for the storage.
func WithFactory(factory model.StorageBuilderFactory) Option {
	return func(s *model.StorageFlag) {
		s.Factory = factory
	}
}

var (
	// factoryOnce is used to initialize the factory once.
	factoryOnce sync.Once
	// factory creates a new StorageBuilder by LocType.
	factory = model.BuildersFactory{}
)

// Storage creates a new storage with the given location, repo path prefix and options.
func Storage(location model.Location, repoPathPrefix string, opts ...Option) model.StorageFlag {
	factoryOnce.Do(func() {
		// Register storage builders.
		factory[rs.LocTypeFilestore] = fs.New
	})
	// create a new storage with the given location, repo path prefix and defaults.
	s := model.StorageFlag{
		Location:       location,
		RepoPathPrefix: repoPathPrefix,
		Logger:         &cmdlog.NopLogger{},
		Factory:        &factory,
	}
	// apply storage options.
	for _, opt := range opts {
		opt(&s)
	}
	return s
}