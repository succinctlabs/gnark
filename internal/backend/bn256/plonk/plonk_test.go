// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plonk_test

import (
	"testing"

	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/plonk"
	mockcommitment "github.com/consensys/gnark/crypto/polynomial/bn256/mock_commitment"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/internal/backend/bn256/cs"
	plonkbn256 "github.com/consensys/gnark/internal/backend/bn256/plonk"
	bn256witness "github.com/consensys/gnark/internal/backend/bn256/witness"
	"github.com/consensys/gnark/internal/backend/circuits"
	curve "github.com/consensys/gurvy/bn256"
)

func TestCircuits(t *testing.T) {
	for name, circuit := range circuits.Circuits {
		t.Run(name, func(t *testing.T) {
			assert := plonk.NewAssert(t)
			pcs, err := frontend.Compile(curve.ID, backend.PLONK, circuit.Circuit)
			assert.NoError(err)
			assert.SolvingSucceeded(pcs, circuit.Good)
			assert.SolvingFailed(pcs, circuit.Bad)
		})
	}
}

// TODO WIP -> once everything is clean move this to backend/plonk in assert
func TestProver(t *testing.T) {

	for name, circuit := range circuits.Circuits {
		// name := "range"
		// circuit := circuits.Circuits[name]

		t.Run(name, func(t *testing.T) {

			assert := plonk.NewAssert(t)
			pcs, err := frontend.Compile(curve.ID, backend.PLONK, circuit.Circuit)
			assert.NoError(err)

			spr := pcs.(*cs.SparseR1CS)

			scheme := mockcommitment.Scheme{}
			wPublic := bn256witness.Witness{}
			wPublic.FromPublicAssignment(circuit.Good)
			publicData := plonkbn256.Setup(spr, &scheme, wPublic)

			// correct proof
			{
				wFull := bn256witness.Witness{}
				wFull.FromFullAssignment(circuit.Good)
				proof := plonkbn256.Prove(spr, publicData, wFull)

				v := plonkbn256.VerifyRaw(proof, publicData, wPublic)

				if !v {
					t.Fatal("Correct proof verification failed")
				}
			}

			//wrong proof
			{
				wFull := bn256witness.Witness{}
				wFull.FromFullAssignment(circuit.Bad)
				proof := plonkbn256.Prove(spr, publicData, wFull)

				v := plonkbn256.VerifyRaw(proof, publicData, wPublic)

				if v {
					t.Fatal("Wrong proof verification should have failed")
				}
			}
		})

	}
}
