// Copyright 2023 The Kanister Authors.
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

package datamover

import (
	"context"

	"github.com/pkg/errors"

	crv1alpha1 "github.com/kanisterio/kanister/pkg/apis/cr/v1alpha1"
	"github.com/kanisterio/kanister/pkg/kopia/snapshot"
	"github.com/kanisterio/kanister/pkg/param"
)

type Profile struct {
	OutputName string
	Profile    *param.Profile
	SnapJSON   string
}

func (p *Profile) Pull(ctx context.Context, sourcePath, destinationPath string) error {
	if p.Profile.Location.Type == crv1alpha1.LocationTypeKopia {
		if p.SnapJSON == "" {
			return errors.New("kopia snapshot information is required to pull data using kopia")
		}
		kopiaSnap, err := snapshot.UnmarshalKopiaSnapshot(p.SnapJSON)
		if err != nil {
			return err
		}
		if err = connectToKopiaServer(ctx, p.Profile); err != nil {
			return err
		}
		return kopiaLocationPull(ctx, kopiaSnap.ID, destinationPath, sourcePath, p.Profile.Credential.KopiaServerSecret.Password)
	}
	target, err := targetWriter(sourcePath)
	if err != nil {
		return err
	}
	return locationPull(ctx, p.Profile, destinationPath, target)
}

func (p *Profile) Push(ctx context.Context, sourcePath, destinationPath string) error {
	if p.Profile.Location.Type == crv1alpha1.LocationTypeKopia {
		if err := connectToKopiaServer(ctx, p.Profile); err != nil {
			return err
		}
		return kopiaLocationPush(ctx, destinationPath, p.OutputName, sourcePath, p.Profile.Credential.KopiaServerSecret.Password)
	}
	source, err := sourceReader(sourcePath)
	if err != nil {
		return err
	}
	return locationPush(ctx, p.Profile, destinationPath, source)
}

func (p *Profile) Delete(ctx context.Context, destinationPath string) error {
	if p.Profile.Location.Type == crv1alpha1.LocationTypeKopia {
		if p.SnapJSON == "" {
			return errors.New("kopia snapshot information is required to delete data using kopia")
		}
		kopiaSnap, err := snapshot.UnmarshalKopiaSnapshot(p.SnapJSON)
		if err != nil {
			return err
		}
		if err = connectToKopiaServer(ctx, p.Profile); err != nil {
			return err
		}
		return kopiaLocationDelete(ctx, kopiaSnap.ID, destinationPath, p.Profile.Credential.KopiaServerSecret.Password)
	}
	return locationDelete(ctx, p.Profile, destinationPath)
}
