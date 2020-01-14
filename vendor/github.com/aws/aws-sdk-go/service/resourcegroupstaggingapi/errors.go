// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package resourcegroupstaggingapi

const (

	// ErrCodeConcurrentModificationException for service response error code
	// "ConcurrentModificationException".
	//
	// The target of the operation is currently being modified by a different request.
	// Try again later.
	ErrCodeConcurrentModificationException = "ConcurrentModificationException"

	// ErrCodeConstraintViolationException for service response error code
	// "ConstraintViolationException".
	//
	// The request was denied because performing this operation violates a constraint.
	//
	// Some of the reasons in the following list might not apply to this specific
	// operation.
	//
	//    * You must meet the prerequisites for using tag policies. For information,
	//    see Prerequisites and Permissions for Using Tag Policies (http://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_tag-policies-prereqs.html)
	//    in the AWS Organizations User Guide.
	//
	//    * You must enable the tag policies service principal (tagpolicies.tag.amazonaws.com)
	//    to integrate with AWS Organizations For information, see EnableAWSServiceAccess
	//    (http://docs.aws.amazon.com/organizations/latest/APIReference/API_EnableAWSServiceAccess.html).
	//
	//    * You must have a tag policy attached to the organization root, an OU,
	//    or an account.
	ErrCodeConstraintViolationException = "ConstraintViolationException"

	// ErrCodeInternalServiceException for service response error code
	// "InternalServiceException".
	//
	// The request processing failed because of an unknown error, exception, or
	// failure. You can retry the request.
	ErrCodeInternalServiceException = "InternalServiceException"

	// ErrCodeInvalidParameterException for service response error code
	// "InvalidParameterException".
	//
	// This error indicates one of the following:
	//
	//    * A parameter is missing.
	//
	//    * A malformed string was supplied for the request parameter.
	//
	//    * An out-of-range value was supplied for the request parameter.
	//
	//    * The target ID is invalid, unsupported, or doesn't exist.
	//
	//    * You can't access the Amazon S3 bucket for report storage. For more information,
	//    see Additional Requirements for Organization-wide Tag Compliance Reports
	//    (http://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_tag-policies-prereqs.html#bucket-policies-org-report)
	//    in the AWS Organizations User Guide.
	ErrCodeInvalidParameterException = "InvalidParameterException"

	// ErrCodePaginationTokenExpiredException for service response error code
	// "PaginationTokenExpiredException".
	//
	// A PaginationToken is valid for a maximum of 15 minutes. Your request was
	// denied because the specified PaginationToken has expired.
	ErrCodePaginationTokenExpiredException = "PaginationTokenExpiredException"

	// ErrCodeThrottledException for service response error code
	// "ThrottledException".
	//
	// The request was denied to limit the frequency of submitted requests.
	ErrCodeThrottledException = "ThrottledException"
)
