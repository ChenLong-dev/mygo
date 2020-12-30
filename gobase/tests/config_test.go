/*
 * @Description:
 * @Author: Chen Long
 * @Date: 2020-12-17 09:33:15
 * @LastEditTime: 2020-12-17 09:37:19
 * @LastEditors: Chen Long
 * @Reference:
 */

package tests

// Basic imports
import (
	"testing"

	"tests/config"

	"github.com/stretchr/testify/suite"
)

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(config.ConfigTestSuite))
}
