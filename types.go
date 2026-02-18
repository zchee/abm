// Copyright 2026 The abm Authors.
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
//
// SPDX-License-Identifier: Apache-2.0

package abm

import (
	"time"
)

// OrgDevicesResponse a response that contains a single organization device resource.
type OrgDevicesResponse struct {
	Data  []OrgDevice        `json:"data,omitzero"`
	Links PagedDocumentLinks `json:"links,omitzero"`
}

// OrgDevice is the data structure that represents an organization device resource.
type OrgDevice struct {
	// The resource’s attributes.
	Attributes *OrgDeviceAttributes `json:"attributes,omitzero"`
	// The opaque resource ID that uniquely identifies the resource.
	ID string `json:"id"`
	// Navigational links that include the self-link.
	Links Links `json:"links,omitzero"`
	// Navigational links to related data and included resource types and IDs.
	Relationships OrgDeviceRelationships `json:"relationships"`
	// The resource type. Always "orgDevices".
	Type string `json:"type,omitzero"`
}

type OrgDeviceAttributesProductFamily string

const (
	ProductFamilyIPhone  OrgDeviceAttributesProductFamily = "iPhone"
	ProductFamilyIPad    OrgDeviceAttributesProductFamily = "iPad"
	ProductFamilyMac     OrgDeviceAttributesProductFamily = "Mac"
	ProductFamilyAppleTV OrgDeviceAttributesProductFamily = "AppleTV"
	ProductFamilyWatch   OrgDeviceAttributesProductFamily = "Watch"
	ProductFamilyVision  OrgDeviceAttributesProductFamily = "Vision"
)

type OrgDeviceAttributesPurchaseSourceType string

const (
	PurchaseSourceTypeApple         OrgDeviceAttributesPurchaseSourceType = "APPLE"
	PurchaseSourceTypeReseller      OrgDeviceAttributesPurchaseSourceType = "RESELLER"
	PurchaseSourceTypeManuallyAdded OrgDeviceAttributesPurchaseSourceType = "MANUALLY_ADDED"
)

type OrgDeviceAttributesStatus string

const (
	StatusAssigned   OrgDeviceAttributesStatus = "ASSIGNED"
	StatusUnAssigned OrgDeviceAttributesStatus = "UNASSIGNED"
)

// OrgDeviceAttributes is the resource’s attributes.
type OrgDeviceAttributes struct {
	// The date and time of adding the device to an organization.
	AddedToOrgDateTime time.Time `json:"addedToOrgDateTime"`
	// The date and time the device was released from an organization.
	// This will be null if the device hasn’t been released.
	// Currently only querying by a single device is supported. Batch device queries aren’t currently supported for this property.
	ReleasedFromOrgDateTime time.Time `json:"releasedFromOrgDateTime,omitzero"`
	// The color of the device.
	Color string `json:"color"`
	// The capacity of the device.
	DeviceCapacity string `json:"deviceCapacity"`
	// The model name.
	DeviceModel string `json:"deviceModel"`
	// The device’s EID (if available).
	EID string `json:"eid,omitzero"`
	// The device’s IMEI (if available).
	IMEI []string `json:"imei,omitempty"`
	// The device’s MEID (if available).
	MEID []string `json:"meid,omitempty"`
	// The device’s Wi-Fi MAC address.
	WifiMacAddress string `json:"wifiMacAddress,omitzero"`
	// The device’s Bluetooth MAC address.
	BluetoothMacAddress string `json:"bluetoothMacAddress,omitempty"`
	// The device’s built-in Ethernet MAC addresses.
	EthernetMacAddress []string `json:"ethernetMacAddress,omitempty"`
	// The date and time of placing the device’s order.
	OrderDateTime time.Time `json:"orderDateTime"`
	// The order number of the device.
	OrderNumber string `json:"orderNumber"`
	// The part number of the device.
	PartNumber string `json:"partNumber"`
	// The device’s Apple product family.
	ProductFamily OrgDeviceAttributesProductFamily `json:"productFamily"`
	// The device’s product type.
	ProductType string `json:"productType"`
	// The device’s purchase source type.
	PurchaseSourceType OrgDeviceAttributesPurchaseSourceType `json:"purchaseSourceType"`
	// The unique ID of the purchase source type: Apple Customer Number or Reseller Number
	PurchaseSourceID string `json:"purchaseSourceId,omitzero"`
	// The device’s serial number.
	SerialNumber string `json:"serialNumber"`
	// The devices status. If [StatusAssigned], use a separate API to get the information of the [StatusUnAssigned] server.
	Status OrgDeviceAttributesStatus `json:"status"`
	// The date and time of the most-recent update for the device.
	UpdatedDateTime time.Time `json:"updatedDateTime"`
}

// Links navigational links.
type Links struct {
	// The link that produces the current document.
	Self string `json:"self"`
}

// The relationships you include in the request, and those that you can operate on.
type OrgDeviceRelationships struct {
	// The relationship representing a device and its assigned device management
	// service.
	AssignedServer *OrgDeviceRelationshipsAssignedServer `json:"assignedServer,omitempty"`

	// The relationship representing a device and its AppleCare Coverage.
	AppleCareCoverage *OrgDeviceRelationshipsAppleCareCoverage `json:"appleCareCoverage,omitempty"`
}

// The links that describe the relationship between the resources.
type OrgDeviceRelationshipsAssignedServer struct {
	// Links corresponds to the JSON schema field "links".
	Links *RelationshipLinks `json:"links,omitempty"`
}

// The links that describe the relationship between the resources.
type OrgDeviceRelationshipsAppleCareCoverage struct {
	// Links corresponds to the JSON schema field "links".
	Links *RelationshipLinks `json:"links,omitempty"`
}

// Links related to the response document, including self-links.
type RelationshipLinks struct {
	// Include corresponds to the JSON schema field "include".
	Include string `json:"include,omitzero"`

	// The link to the related documents.
	Related string `json:"related,omitzero"`

	// The link that produces the current document.
	Self string `json:"self,omitzero"`
}

// PagedDocumentLinks navigational links that include the self-link.
type PagedDocumentLinks struct {
	Links
	// The link to the first page of documents.
	First string `json:"first,omitzero"`
	// The link to the next page of documents.
	Next string `json:"next,omitzero"`
}
