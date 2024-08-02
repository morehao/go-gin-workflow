package flow

import (
	"errors"
	"fmt"
	"go-gin-workflow/internal/workflow/object/objFlow"

	"github.com/morehao/go-tools/gutils"
)

const (
	NodeTypeStart     objFlow.NodeType = "start"     // 开始，工作流的起点
	NodeTypeRoute     objFlow.NodeType = "route"     // 路由，根据条件判断，决定下一步走向
	NodeTypeCondition objFlow.NodeType = "condition" // 条件，用于判断，决定下一步走向
	NodeTypeApprover  objFlow.NodeType = "approver"  // 审批人
	NodeTypeNotifier  objFlow.NodeType = "notifier"  // 通知人
)

var NodeTypeList = []objFlow.NodeType{NodeTypeStart, NodeTypeRoute, NodeTypeCondition, NodeTypeApprover, NodeTypeNotifier}

// IsValidProcessConfig 检查流程配置是否有效
func IsValidProcessConfig(node *objFlow.Node) error {
	// 节点名称是否有效
	if len(node.NodeID) == 0 {
		return errors.New("节点的【nodeId】不能为空！！")
	}
	// 检查类型是否有效
	if len(node.Type) == 0 {
		return fmt.Errorf("节点【%s】的类型不能为空！！", node.NodeID)
	}
	var flag = false
	for _, nodeType := range NodeTypeList {
		if nodeType == node.Type {
			flag = true
			break
		}
	}
	if !flag {
		return fmt.Errorf("节点【%s】的类型为【%s】，为无效类型, 有效类型为%s", node.NodeID, node.Type, gutils.ToJsonString(NodeTypeList))
	}
	// 当前节点是否设置有审批人
	if node.Type == NodeTypeApprover || node.Type == NodeTypeNotifier {
		if node.Properties == nil || node.Properties.ActionerRules == nil {
			return fmt.Errorf("节点【%s】的Properties属性不能为空，如：`\"properties\": {\"actionerRules\": [{\"type\": \"target_label\",\"labelNames\": \"人事\",\"memberCount\": 1,\"actType\": \"and\"}],}`", node.NodeID)
		}
	}
	// 条件节点是否存在
	if node.ConditionNodes != nil { // 存在条件节点
		if len(node.ConditionNodes) == 1 {
			return fmt.Errorf("节点【%s】条件节点下的节点数必须大于1", node.NodeID)
		}
		// 根据条件变量选择节点索引
		for _, conditionNode := range node.ConditionNodes {
			if conditionNode.Properties == nil {
				return fmt.Errorf("节点【%s】的Properties对象为空值！！", conditionNode.NodeID)
			}
			if len(conditionNode.Properties.Conditions) == 0 {
				return fmt.Errorf("节点【%s】的Conditions对象为空值！！", conditionNode.NodeID)
			}
			err := IsValidProcessConfig(conditionNode)
			if err != nil {
				return err
			}
		}
	}

	// 子节点是否存在
	if node.ChildNode != nil {
		return IsValidProcessConfig(node.ChildNode)
	}
	return nil
}
