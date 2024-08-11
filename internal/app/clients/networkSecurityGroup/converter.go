package networkSecurityGroup

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"strconv"
	"x/internal/app/model"
)

func convert(rule model.SecurityRule) *armnetwork.SecurityRule {
	return &armnetwork.SecurityRule{
		Name: to.Ptr(rule.Name),
		Properties: &armnetwork.SecurityRulePropertiesFormat{
			SourceAddressPrefix:      to.Ptr(rule.SourceAddressPrefix),
			SourcePortRange:          to.Ptr(rule.SourcePortRange),
			DestinationAddressPrefix: to.Ptr(rule.DestinationAddressPrefix),
			DestinationPortRange:     to.Ptr(rule.DestinationPortRange),
			Protocol:                 to.Ptr(fetchProtocolType(rule.Protocol)),
			Access:                   to.Ptr(fetchAccessType(rule.Access)),
			Priority:                 to.Ptr[int32](fetchPriorityType(rule.Priority)),
			Description:              to.Ptr(rule.Description),
			Direction:                to.Ptr(fetchDirectionType(rule.Direction)),
		},
	}
}

func fetchProtocolType(from string) armnetwork.SecurityRuleProtocol {
	switch from {
	case "tcp":
		return armnetwork.SecurityRuleProtocolTCP
	case "udp":
		return armnetwork.SecurityRuleProtocolUDP
	case "icmp":
		return armnetwork.SecurityRuleProtocolIcmp
	case "esp":
		return armnetwork.SecurityRuleProtocolEsp
	case "ah":
		return armnetwork.SecurityRuleProtocolAh
	case "*":
		return armnetwork.SecurityRuleProtocolAsterisk
	default:
		return ""
	}
}

func fetchAccessType(from string) armnetwork.SecurityRuleAccess {
	switch from {
	case "allow":
		return armnetwork.SecurityRuleAccessAllow
	case "deny":
		return armnetwork.SecurityRuleAccessDeny
	default:
		return ""
	}
}

func fetchDirectionType(from string) armnetwork.SecurityRuleDirection {
	switch from {
	case "inbound":
		return armnetwork.SecurityRuleDirectionInbound
	case "outbound":
		return armnetwork.SecurityRuleDirectionOutbound
	default:
		return ""
	}
}

func fetchPriorityType(from string) int32 {
	res, _ := strconv.ParseInt(from, 10, 32)
	return int32(res)
}
