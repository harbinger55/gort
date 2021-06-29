/*
 * Copyright 2021 The Gort Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/getgort/gort/client"
)

const (
	groupRemoveRoleUse   = "remove-role"
	groupRemoveRoleShort = "Remove a role from an existing group"
	groupRemoveRoleLong  = "Remove a role from an existing group."
)

// GetGroupRemoveRoleCmd is a command
func GetGroupRemoveRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   groupRemoveRoleUse,
		Short: groupRemoveRoleShort,
		Long:  groupRemoveRoleLong,
		RunE:  groupRemoveRoleCmd,
		Args:  cobra.ExactArgs(2),
	}

	return cmd
}

func groupRemoveRoleCmd(cmd *cobra.Command, args []string) error {
	groupname := args[0]
	rolename := args[1]

	gortClient, err := client.Connect(FlagGortProfile)
	if err != nil {
		return err
	}

	err = gortClient.GroupRoleDelete(groupname, rolename)
	if err != nil {
		return err
	}

	fmt.Printf("Role removed from %s: %s\n", groupname, rolename)

	return nil
}
