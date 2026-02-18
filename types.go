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

import "time"

// OrgDevicesResponse contains a list of organization device resources.
type OrgDevicesResponse struct {
	Data  []OrgDevice        `json:"data"`
	Links PagedDocumentLinks `json:"links"`
	Meta  *PagingInformation `json:"meta,omitempty"`
}

// OrgDeviceResponse contains a single organization device resource.
type OrgDeviceResponse struct {
	Data  OrgDevice     `json:"data"`
	Links DocumentLinks `json:"links"`
}

// OrgDevice represents an organization device resource.
type OrgDevice struct {
	Attributes    *OrgDeviceAttributes    `json:"attributes,omitempty"`
	ID            string                  `json:"id"`
	Links         *ResourceLinks          `json:"links,omitempty"`
	Relationships *OrgDeviceRelationships `json:"relationships,omitempty"`
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
	AddedToOrgDateTime      *time.Time                            `json:"addedToOrgDateTime,omitempty"`
	ReleasedFromOrgDateTime *time.Time                            `json:"releasedFromOrgDateTime,omitempty"`
	Color                   string                                `json:"color,omitempty"`
	DeviceCapacity          string                                `json:"deviceCapacity,omitempty"`
	DeviceModel             string                                `json:"deviceModel,omitempty"`
	EID                     string                                `json:"eid,omitempty"`
	IMEI                    []string                              `json:"imei,omitempty"`
	MEID                    []string                              `json:"meid,omitempty"`
	WifiMacAddress          []string                              `json:"wifiMacAddress,omitempty"`
	BluetoothMacAddress     []string                              `json:"bluetoothMacAddress,omitempty"`
	EthernetMacAddress      []string                              `json:"ethernetMacAddress,omitempty"`
	OrderDateTime           *time.Time                            `json:"orderDateTime,omitempty"`
	OrderNumber             string                                `json:"orderNumber,omitempty"`
	PartNumber              string                                `json:"partNumber,omitempty"`
	ProductFamily           OrgDeviceAttributesProductFamily      `json:"productFamily,omitempty"`
	ProductType             string                                `json:"productType,omitempty"`
	PurchaseSourceType      OrgDeviceAttributesPurchaseSourceType `json:"purchaseSourceType,omitempty"`
	PurchaseSourceID        string                                `json:"purchaseSourceId,omitempty"`
	SerialNumber            string                                `json:"serialNumber,omitempty"`
	Status                  OrgDeviceAttributesStatus             `json:"status,omitempty"`
	UpdatedDateTime         *time.Time                            `json:"updatedDateTime,omitempty"`
}

// OrgDeviceRelationships contains links to relationship resources for an org device.
type OrgDeviceRelationships struct {
	AssignedServer    *OrgDeviceRelationshipsAssignedServer    `json:"assignedServer,omitempty"`
	AppleCareCoverage *OrgDeviceRelationshipsAppleCareCoverage `json:"appleCareCoverage,omitempty"`
}

// OrgDeviceRelationshipsAssignedServer describes assigned-server relationship links.
type OrgDeviceRelationshipsAssignedServer struct {
	Links *RelationshipLinks `json:"links,omitempty"`
}

// OrgDeviceRelationshipsAppleCareCoverage describes apple-care relationship links.
type OrgDeviceRelationshipsAppleCareCoverage struct {
	Links *RelationshipLinks `json:"links,omitempty"`
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

// MdmServersResponse contains a list of MDM server resources.
type MdmServersResponse struct {
	Data     []MdmServer        `json:"data"`
	Included []OrgDevice        `json:"included,omitempty"`
	Links    PagedDocumentLinks `json:"links"`
	Meta     *PagingInformation `json:"meta,omitempty"`
}

// MdmServerResponse contains a single MDM server resource.
type MdmServerResponse struct {
	Data     MdmServer     `json:"data"`
	Included []OrgDevice   `json:"included,omitempty"`
	Links    DocumentLinks `json:"links"`
}

// MdmServer is a device management service resource.
type MdmServer struct {
	Attributes    *MdmServerAttributes    `json:"attributes,omitempty"`
	ID            string                  `json:"id"`
	Relationships *MdmServerRelationships `json:"relationships,omitempty"`
	Type          string                  `json:"type"`
}

// MdmServerAttributes are fields describing an MDM server.
type MdmServerAttributes struct {
	CreatedDateTime *time.Time `json:"createdDateTime,omitempty"`
	ServerName      string     `json:"serverName,omitempty"`
	ServerType      string     `json:"serverType,omitempty"`
	UpdatedDateTime *time.Time `json:"updatedDateTime,omitempty"`
}

// MdmServerRelationships contains relationship resources for an MDM server.
type MdmServerRelationships struct {
	Devices *MdmServerRelationshipsDevices `json:"devices,omitempty"`
}

// MdmServerRelationshipsDevices represents the devices relationship in an MDM server.
type MdmServerRelationshipsDevices struct {
	Data  []MdmServerRelationshipsDevicesData `json:"data,omitempty"`
	Links *RelationshipLinks                  `json:"links,omitempty"`
	Meta  *PagingInformation                  `json:"meta,omitempty"`
}

// MdmServerRelationshipsDevicesData is an org-device linkage in an MDM-server relationship.
type MdmServerRelationshipsDevicesData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// MdmServerDevicesLinkagesResponse contains org-device linkages for a specific MDM server.
type MdmServerDevicesLinkagesResponse struct {
	Data  []MdmServerDevicesLinkageData `json:"data"`
	Links PagedDocumentLinks            `json:"links"`
	Meta  *PagingInformation            `json:"meta,omitempty"`
}

// MdmServerDevicesLinkageData contains an org-device linkage entry.
type MdmServerDevicesLinkageData struct {
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
	Attributes *OrgDeviceActivityAttributes `json:"attributes,omitempty"`
	ID         string                       `json:"id"`
	Links      *ResourceLinks               `json:"links,omitempty"`
	Type       string                       `json:"type"`
}

// OrgDeviceActivityAttributes are fields describing an org-device activity.
type OrgDeviceActivityAttributes struct {
	CompletedDateTime *time.Time `json:"completedDateTime,omitempty"`
	CreatedDateTime   *time.Time `json:"createdDateTime,omitempty"`
	DownloadURL       string     `json:"downloadUrl,omitempty"`
	Status            string     `json:"status,omitempty"`
	SubStatus         string     `json:"subStatus,omitempty"`
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
	MdmServer OrgDeviceActivityCreateRequestDataRelationshipsMdmServer `json:"mdmServer"`
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

// OrgDeviceActivityCreateRequestDataRelationshipsMdmServer links an MDM server in activity creation.
type OrgDeviceActivityCreateRequestDataRelationshipsMdmServer struct {
	Data OrgDeviceActivityCreateRequestDataRelationshipsMdmServerData `json:"data"`
}

// OrgDeviceActivityCreateRequestDataRelationshipsMdmServerData is an MDM-server linkage used in activity creation.
type OrgDeviceActivityCreateRequestDataRelationshipsMdmServerData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// AppleCareCoverageResponse contains AppleCare coverage resources for a device.
type AppleCareCoverageResponse struct {
	Data  []AppleCareCoverage `json:"data"`
	Links PagedDocumentLinks  `json:"links"`
	Meta  *PagingInformation  `json:"meta,omitempty"`
}

// AppleCareCoverage contains AppleCare coverage data.
type AppleCareCoverage struct {
	Attributes *AppleCareCoverageAttributes `json:"attributes,omitempty"`
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
	AgreementNumber        string                       `json:"agreementNumber,omitempty"`
	ContractCancelDateTime *time.Time                   `json:"contractCancelDateTime,omitempty"`
	Description            string                       `json:"description,omitempty"`
	EndDateTime            *time.Time                   `json:"endDateTime,omitempty"`
	IsCanceled             *bool                        `json:"isCanceled,omitempty"`
	IsRenewable            *bool                        `json:"isRenewable,omitempty"`
	PaymentType            AppleCareCoveragePaymentType `json:"paymentType,omitempty"`
	StartDateTime          *time.Time                   `json:"startDateTime,omitempty"`
	Status                 AppleCareCoverageStatus      `json:"status,omitempty"`
}

// DocumentLinks contains links related to the current document.
type DocumentLinks struct {
	Self string `json:"self"`
}

// ResourceLinks contains self-links for a requested resource.
type ResourceLinks struct {
	Self string `json:"self,omitempty"`
}

// RelationshipLinks contains links for a relationship block.
type RelationshipLinks struct {
	Include string `json:"include,omitempty"`
	Related string `json:"related,omitempty"`
	Self    string `json:"self,omitempty"`
}

// PagedDocumentLinks contains navigation links for paginated responses.
type PagedDocumentLinks struct {
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Self  string `json:"self"`
}

// PagingInformation contains pagination metadata.
type PagingInformation struct {
	Paging PagingInformationPaging `json:"paging"`
}

// PagingInformationPaging contains pagination state values.
type PagingInformationPaging struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"nextCursor,omitempty"`
	Total      *int   `json:"total,omitempty"`
}

// ErrorResponse contains ABM API errors.
type ErrorResponse struct {
	Errors []ErrorResponseError `json:"errors,omitempty"`
}

// ErrorResponseError contains one ABM API error object.
type ErrorResponseError struct {
	Code   string         `json:"code"`
	Detail string         `json:"detail"`
	ID     string         `json:"id,omitempty"`
	Links  *ErrorLinks    `json:"links,omitempty"`
	Meta   map[string]any `json:"meta,omitempty"`
	Source *ErrorSource   `json:"source,omitempty"`
	Status string         `json:"status"`
	Title  string         `json:"title"`
}

// ErrorLinks contains links related to an error object.
type ErrorLinks struct {
	About      string `json:"about,omitempty"`
	Associated any    `json:"associated,omitempty"`
}

// ErrorLinksAssociated contains structured associated error-link details.
type ErrorLinksAssociated struct {
	Href string                    `json:"href,omitempty"`
	Meta *ErrorLinksAssociatedMeta `json:"meta,omitempty"`
}

// ErrorLinksAssociatedMeta contains metadata for associated error links.
type ErrorLinksAssociatedMeta struct {
	Source string `json:"source,omitempty"`
}

// ErrorSource describes JSON Pointer or parameter source context for an error.
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// JSONPointer contains a JSON Pointer source.
type JSONPointer struct {
	Pointer string `json:"pointer"`
}

// Parameter contains an HTTP parameter source.
type Parameter struct {
	Parameter string `json:"parameter"`
}
