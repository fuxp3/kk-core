/*
 Copyright 2021 The KubeSphere Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package module

import (
	"kk-core/core/cache"
	"kk-core/core/connector"
	"kk-core/core/ending"
)

type Module interface {
	IsSkip() bool
	Default(runtime connector.Runtime, pipelineCache *cache.Cache, moduleCache *cache.Cache)
	Init()
	Is() string
	Run(result *ending.ModuleResult)
	Until() (*bool, error)
	Slogan()
	AutoAssert()
	AppendPostHook(h PostHookInterface)
	CallPostHook(result *ending.ModuleResult) error
}
