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

// OrgDevicesResponse contains a list of organization device resources.
type OrgDevicesResponse struct {
	Data  []OrgDevice        `json:"data"`
	Links PagedDocumentLinks `json:"links"`
	Meta  *PagingInformation `json:"meta,omitzero"`
}

// OrgDeviceResponse contains a single organization device resource.
type OrgDeviceResponse struct {
	Data  OrgDevice     `json:"data"`
	Links DocumentLinks `json:"links"`
}

// OrgDevice represents an organization device resource.
type OrgDevice struct {
	Attributes    *OrgDeviceAttributes    `json:"attributes,omitzero"`
	ID            string                  `json:"id"`
	Links         *ResourceLinks          `json:"links,omitzero"`
	Relationships *OrgDeviceRelationships `json:"relationships,omitzero"`
	Type          string                  `json:"type"`
}

// OrgDeviceAttributesProductFamily is the product family of an organization device.
type OrgDeviceAttributesProductFamily string

const (
	ProductFamilyIPhone  OrgDeviceAttributesProductFamily = "iPhone"
	ProductFamilyIPad    OrgDeviceAttributesProductFamily = "iPad"
	ProductFamilyMac     OrgDeviceAttributesProductFamily = "Mac"
	ProductFamilyAppleTV OrgDeviceAttributesProductFamily = "AppleTV"
	ProductFamilyWatch   OrgDeviceAttributesProductFamily = "Watch"
	ProductFamilyVision  OrgDeviceAttributesProductFamily = "Vision"
)

// OrgDeviceAttributesPurchaseSourceType is the purchase source type of an organization device.
type OrgDeviceAttributesPurchaseSourceType string

const (
	PurchaseSourceTypeApple         OrgDeviceAttributesPurchaseSourceType = "APPLE"
	PurchaseSourceTypeReseller      OrgDeviceAttributesPurchaseSourceType = "RESELLER"
	PurchaseSourceTypeManuallyAdded OrgDeviceAttributesPurchaseSourceType = "MANUALLY_ADDED"
)

// OrgDeviceAttributesStatus is the assignment status of an organization device.
type OrgDeviceAttributesStatus string

const (
	StatusAssigned   OrgDeviceAttributesStatus = "ASSIGNED"
	StatusUnAssigned OrgDeviceAttributesStatus = "UNASSIGNED"
)

// OrgDeviceAttributes contains attributes for an organization device resource.
type OrgDeviceAttributes struct {
	AddedToOrgDateTime      time.Time                             `json:"addedToOrgDateTime,omitzero"`
	ReleasedFromOrgDateTime time.Time                             `json:"releasedFromOrgDateTime,omitzero"`
	Color                   string                                `json:"color,omitzero"`
	DeviceCapacity          string                                `json:"deviceCapacity,omitzero"`
	DeviceModel             string                                `json:"deviceModel,omitzero"`
	EID                     string                                `json:"eid,omitzero"`
	IMEI                    []string                              `json:"imei,omitempty"`
	MEID                    []string                              `json:"meid,omitempty"`
	WifiMacAddress          []string                              `json:"wifiMacAddress,omitempty"`
	BluetoothMacAddress     []string                              `json:"bluetoothMacAddress,omitempty"`
	EthernetMacAddress      []string                              `json:"ethernetMacAddress,omitempty"`
	OrderDateTime           time.Time                             `json:"orderDateTime,omitzero"`
	OrderNumber             string                                `json:"orderNumber,omitzero"`
	PartNumber              string                                `json:"partNumber,omitzero"`
	ProductFamily           OrgDeviceAttributesProductFamily      `json:"productFamily,omitzero"`
	ProductType             string                                `json:"productType,omitzero"`
	PurchaseSourceType      OrgDeviceAttributesPurchaseSourceType `json:"purchaseSourceType,omitzero"`
	PurchaseSourceID        string                                `json:"purchaseSourceId,omitzero"`
	SerialNumber            string                                `json:"serialNumber,omitzero"`
	Status                  OrgDeviceAttributesStatus             `json:"status,omitzero"`
	UpdatedDateTime         time.Time                             `json:"updatedDateTime,omitzero"`
}

// OrgDeviceRelationships contains links to relationship resources for an org device.
type OrgDeviceRelationships struct {
	AssignedServer    *OrgDeviceRelationshipsAssignedServer    `json:"assignedServer,omitzero"`
	AppleCareCoverage *OrgDeviceRelationshipsAppleCareCoverage `json:"appleCareCoverage,omitzero"`
}

// OrgDeviceRelationshipsAssignedServer describes assigned-server relationship links.
type OrgDeviceRelationshipsAssignedServer struct {
	Links *RelationshipLinks `json:"links,omitzero"`
}

// OrgDeviceRelationshipsAppleCareCoverage describes apple-care relationship links.
type OrgDeviceRelationshipsAppleCareCoverage struct {
	Links *RelationshipLinks `json:"links,omitzero"`
}

// OrgDeviceAssignedServerLinkageResponse contains the assigned server linkage for a device.
type OrgDeviceAssignedServerLinkageResponse struct {
	Data  OrgDeviceAssignedServerLinkageData `json:"data"`
	Links DocumentLinks                      `json:"links"`
}

// OrgDeviceAssignedServerLinkageData is the assigned server linkage object.
type OrgDeviceAssignedServerLinkageData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// MDMServersResponse contains a list of MDM server resources.
type MDMServersResponse struct {
	Data     []MDMServer        `json:"data"`
	Included []OrgDevice        `json:"included,omitempty"`
	Links    PagedDocumentLinks `json:"links"`
	Meta     *PagingInformation `json:"meta,omitzero"`
}

// MDMServerResponse contains a single MDM server resource.
type MDMServerResponse struct {
	Data     MDMServer     `json:"data"`
	Included []OrgDevice   `json:"included,omitempty"`
	Links    DocumentLinks `json:"links"`
}

// MDMServer is a device management service resource.
type MDMServer struct {
	Attributes    *MDMServerAttributes    `json:"attributes,omitzero"`
	ID            string                  `json:"id"`
	Relationships *MDMServerRelationships `json:"relationships,omitzero"`
	Type          string                  `json:"type"`
}

// MDMServerAttributes are fields describing an MDM server.
type MDMServerAttributes struct {
	CreatedDateTime time.Time `json:"createdDateTime,omitzero"`
	ServerName      string    `json:"serverName,omitzero"`
	ServerType      string    `json:"serverType,omitzero"`
	UpdatedDateTime time.Time `json:"updatedDateTime,omitzero"`
}

// MDMServerRelationships contains relationship resources for an MDM server.
type MDMServerRelationships struct {
	Devices *MDMServerRelationshipsDevices `json:"devices,omitzero"`
}

// MDMServerRelationshipsDevices represents the devices relationship in an MDM server.
type MDMServerRelationshipsDevices struct {
	Data  []MDMServerRelationshipsDevicesData `json:"data,omitempty"`
	Links *RelationshipLinks                  `json:"links,omitzero"`
	Meta  *PagingInformation                  `json:"meta,omitzero"`
}

// MDMServerRelationshipsDevicesData is an org-device linkage in an MDM-server relationship.
type MDMServerRelationshipsDevicesData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// MDMServerDevicesLinkagesResponse contains org-device linkages for a specific MDM server.
type MDMServerDevicesLinkagesResponse struct {
	Data  []MDMServerDevicesLinkageData `json:"data"`
	Links PagedDocumentLinks            `json:"links"`
	Meta  *PagingInformation            `json:"meta,omitzero"`
}

// MDMServerDevicesLinkageData contains an org-device linkage entry.
type MDMServerDevicesLinkageData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// OrgDeviceActivityResponse contains a single org-device activity resource.
type OrgDeviceActivityResponse struct {
	Data  OrgDeviceActivity `json:"data"`
	Links DocumentLinks     `json:"links"`
}

// OrgDeviceActivity is an activity resource for assigning or unassigning devices.
type OrgDeviceActivity struct {
	Attributes *OrgDeviceActivityAttributes `json:"attributes,omitzero"`
	ID         string                       `json:"id"`
	Links      *ResourceLinks               `json:"links,omitzero"`
	Type       string                       `json:"type"`
}

// OrgDeviceActivityAttributes are fields describing an org-device activity.
type OrgDeviceActivityAttributes struct {
	CompletedDateTime time.Time `json:"completedDateTime,omitzero"`
	CreatedDateTime   time.Time `json:"createdDateTime,omitzero"`
	DownloadURL       string    `json:"downloadUrl,omitzero"`
	Status            string    `json:"status,omitzero"`
	SubStatus         string    `json:"subStatus,omitzero"`
}

// OrgDeviceActivityType is the type of an org-device activity.
type OrgDeviceActivityType string

const (
	OrgDeviceActivityTypeAssignDevices   OrgDeviceActivityType = "ASSIGN_DEVICES"
	OrgDeviceActivityTypeUnassignDevices OrgDeviceActivityType = "UNASSIGN_DEVICES"
)

// OrgDeviceActivityCreateRequest is the request payload for creating org-device activities.
type OrgDeviceActivityCreateRequest struct {
	Data OrgDeviceActivityCreateRequestData `json:"data"`
}

// OrgDeviceActivityCreateRequestData is the data section of activity creation requests.
type OrgDeviceActivityCreateRequestData struct {
	Attributes    OrgDeviceActivityCreateRequestDataAttributes    `json:"attributes"`
	Relationships OrgDeviceActivityCreateRequestDataRelationships `json:"relationships"`
	Type          string                                          `json:"type"`
}

// OrgDeviceActivityCreateRequestDataAttributes are activity creation attributes.
type OrgDeviceActivityCreateRequestDataAttributes struct {
	ActivityType OrgDeviceActivityType `json:"activityType"`
}

// OrgDeviceActivityCreateRequestDataRelationships are activity creation relationships.
type OrgDeviceActivityCreateRequestDataRelationships struct {
	Devices   OrgDeviceActivityCreateRequestDataRelationshipsDevices   `json:"devices"`
	MDMServer OrgDeviceActivityCreateRequestDataRelationshipsMDMServer `json:"mdmServer"`
}

// OrgDeviceActivityCreateRequestDataRelationshipsDevices links devices in activity creation.
type OrgDeviceActivityCreateRequestDataRelationshipsDevices struct {
	Data []OrgDeviceActivityCreateRequestDataRelationshipsDevicesData `json:"data"`
}

// OrgDeviceActivityCreateRequestDataRelationshipsDevicesData is a device linkage used in activity creation.
type OrgDeviceActivityCreateRequestDataRelationshipsDevicesData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// OrgDeviceActivityCreateRequestDataRelationshipsMDMServer links an MDM server in activity creation.
type OrgDeviceActivityCreateRequestDataRelationshipsMDMServer struct {
	Data OrgDeviceActivityCreateRequestDataRelationshipsMDMServerData `json:"data"`
}

// OrgDeviceActivityCreateRequestDataRelationshipsMDMServerData is an MDM-server linkage used in activity creation.
type OrgDeviceActivityCreateRequestDataRelationshipsMDMServerData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// AppleCareCoverageResponse contains AppleCare coverage resources for a device.
type AppleCareCoverageResponse struct {
	Data  []AppleCareCoverage `json:"data"`
	Links PagedDocumentLinks  `json:"links"`
	Meta  *PagingInformation  `json:"meta,omitzero"`
}

// AppleCareCoverage contains AppleCare coverage data.
type AppleCareCoverage struct {
	Attributes *AppleCareCoverageAttributes `json:"attributes,omitzero"`
	ID         string                       `json:"id"`
	Type       string                       `json:"type"`
}

// AppleCareCoveragePaymentType is the payment type of an AppleCare coverage plan.
type AppleCareCoveragePaymentType string

const (
	AppleCareCoveragePaymentTypeABESubscription AppleCareCoveragePaymentType = "ABE_SUBSCRIPTION"
	AppleCareCoveragePaymentTypePaidUpFront     AppleCareCoveragePaymentType = "PAID_UP_FRONT"
	AppleCareCoveragePaymentTypeSubscription    AppleCareCoveragePaymentType = "SUBSCRIPTION"
	AppleCareCoveragePaymentTypeNone            AppleCareCoveragePaymentType = "NONE"
)

// AppleCareCoverageStatus is the status of an AppleCare coverage plan.
type AppleCareCoverageStatus string

const (
	AppleCareCoverageStatusActive   AppleCareCoverageStatus = "ACTIVE"
	AppleCareCoverageStatusInactive AppleCareCoverageStatus = "INACTIVE"
)

// AppleCareCoverageAttributes contains AppleCare coverage attributes.
type AppleCareCoverageAttributes struct {
	AgreementNumber        string                       `json:"agreementNumber,omitzero"`
	ContractCancelDateTime time.Time                    `json:"contractCancelDateTime,omitzero"`
	Description            string                       `json:"description,omitzero"`
	EndDateTime            time.Time                    `json:"endDateTime,omitzero"`
	IsCanceled             bool                         `json:"isCanceled,omitzero"`
	IsRenewable            bool                         `json:"isRenewable,omitzero"`
	PaymentType            AppleCareCoveragePaymentType `json:"paymentType,omitzero"`
	StartDateTime          time.Time                    `json:"startDateTime,omitzero"`
	Status                 AppleCareCoverageStatus      `json:"status,omitzero"`
}

// DocumentLinks contains links related to the current document.
type DocumentLinks struct {
	Self string `json:"self"`
}

// ResourceLinks contains self-links for a requested resource.
type ResourceLinks struct {
	Self string `json:"self,omitzero"`
}

// RelationshipLinks contains links for a relationship block.
type RelationshipLinks struct {
	Include string `json:"include,omitzero"`
	Related string `json:"related,omitzero"`
	Self    string `json:"self,omitzero"`
}

// PagedDocumentLinks contains navigation links for paginated responses.
type PagedDocumentLinks struct {
	First string `json:"first,omitzero"`
	Next  string `json:"next,omitzero"`
	Self  string `json:"self"`
}

// PagingInformation contains pagination metadata.
type PagingInformation struct {
	Paging PagingInformationPaging `json:"paging"`
}

// PagingInformationPaging contains pagination state values.
type PagingInformationPaging struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"nextCursor,omitzero"`
	Total      int    `json:"total,omitzero"`
}

// ErrorResponse contains ABM API errors.
type ErrorResponse struct {
	Errors []ErrorResponseError `json:"errors,omitempty"`
}

// ErrorResponseError contains one ABM API error object.
type ErrorResponseError struct {
	Code   string         `json:"code"`
	Detail string         `json:"detail"`
	ID     string         `json:"id,omitzero"`
	Links  *ErrorLinks    `json:"links,omitzero"`
	Meta   map[string]any `json:"meta,omitempty"`
	Source *ErrorSource   `json:"source,omitzero"`
	Status string         `json:"status"`
	Title  string         `json:"title"`
}

// ErrorLinks contains links related to an error object.
type ErrorLinks struct {
	About      string `json:"about,omitzero"`
	Associated any    `json:"associated,omitzero"`
}

// ErrorLinksAssociated contains structured associated error-link details.
type ErrorLinksAssociated struct {
	Href string                    `json:"href,omitzero"`
	Meta *ErrorLinksAssociatedMeta `json:"meta,omitzero"`
}

// ErrorLinksAssociatedMeta contains metadata for associated error links.
type ErrorLinksAssociatedMeta struct {
	Source string `json:"source,omitzero"`
}

// ErrorSource describes JSON Pointer or parameter source context for an error.
type ErrorSource struct {
	Pointer   string `json:"pointer,omitzero"`
	Parameter string `json:"parameter,omitzero"`
}

// JSONPointer contains a JSON Pointer source.
type JSONPointer struct {
	Pointer string `json:"pointer"`
}

// Parameter contains an HTTP parameter source.
type Parameter struct {
	Parameter string `json:"parameter"`
}
