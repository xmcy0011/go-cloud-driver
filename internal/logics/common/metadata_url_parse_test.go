package common

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestIsValidUrl(t *testing.T) {
	Convey("TestIsValidUrl", t, func() {
		Convey("prefix", func() {
			assert.False(t, IsValidUrl(""))
			assert.False(t, IsValidUrl("file://"))
			assert.False(t, IsValidUrl("url:"))
			assert.False(t, IsValidUrl("url:/"))
		})

		Convey("badObjectId", func() {
			assert.False(t, IsValidUrl("url://Abb"))
			assert.False(t, IsValidUrl("url://01HMQDV23J3HMV4CBH2MHSRN45/"))
			assert.False(t, IsValidUrl("url://01HMQDV23J3HMV4CBH2MHSRN45/ab"))
		})

		Convey("valid", func() {
			assert.True(t, IsValidUrl("url://"))
			assert.True(t, IsValidUrl("url://01HMQDV23J3HMV4CBH2MHSRN45"))
			assert.True(t, IsValidUrl("url://01HMQFMBEQKDJP0TTGW9NRFV3V"))
		})
	})
}

func TestGetUrlNodes(t *testing.T) {
	Convey("TestGetUrlNodes", t, func() {
		Convey("valid", func() {
			nodes, err := GetUrlNodes("url://01HMQDV23J3HMV4CBH2MHSRN45")
			assert.NoError(t, err)
			assert.Equal(t, 1, len(nodes))

			nodes, err = GetUrlNodes("url://01HMQDV23J3HMV4CBH2MHSRN45/01HMQGFZ649R7F79ZR76E2S589/01HMQGG3VGDFHSWECAP61CF8PZ")
			assert.NoError(t, err)
			assert.Equal(t, 3, len(nodes))
		})
	})
}

func TestGetUrlLeafNode(t *testing.T) {
	Convey("TestGetUrlLeafNode", t, func() {
		Convey("invalidUrl", func() {
			_, err := GetUrlLeafNode("url:")
			assert.Error(t, ErrInvalidUrlFormat, err)
			_, err = GetUrlLeafNode("file://")
			assert.Error(t, ErrInvalidUrlFormat, err)
			_, err = GetUrlLeafNode("url://01HMQDV23J3HMV4CBH2MHSRN45/abc")
			assert.Error(t, ErrInvalidUrlFormat, err)
			_, err = GetUrlLeafNode("url://01HMQDV23J3HMV4CBH2MHSRN45/01HMQGG3VGDFHSWECAP61CF8PZ/")
			assert.Error(t, ErrInvalidUrlFormat, err)
		})

		Convey("empty", func() {
			_, err := GetUrlLeafNode("url://")
			assert.Error(t, ErrEmptyLeafNode, err)
		})

		Convey("valid", func() {
			node, err := GetUrlLeafNode("url://01HMQDV23J3HMV4CBH2MHSRN45")
			assert.NoError(t, err)
			assert.Equal(t, "01HMQDV23J3HMV4CBH2MHSRN45", node)

			node, err = GetUrlLeafNode("url://01HMQDV23J3HMV4CBH2MHSRN45/01HMQGFZ649R7F79ZR76E2S589/01HMQGG3VGDFHSWECAP61CF8PZ")
			assert.NoError(t, err)
			assert.Equal(t, "01HMQGG3VGDFHSWECAP61CF8PZ", node)
		})
	})
}
