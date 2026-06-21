package routes

import (
	"github.com/gin-gonic/gin"

	"nas-dashboard/internal/api"
)

// registerNetwork 网络接口、Wi-Fi、DNS、代理、PPPoE 等路由。
func registerNetwork(g *gin.RouterGroup) {
	network := g.Group("/network")
	requireAuth(network)
	{
		// 网络接口
		network.GET("/interfaces", api.GetNetworkInterfaces)
		network.GET("/interfaces/ethernet", api.GetEthernetInterfaces)
		network.GET("/interfaces/wifi", api.GetWiFiInterfaces)

		// 接口配置管理（参数化路由）
		network.GET("/interface/:interface/config", api.GetInterfaceConfig)
		network.PUT("/interface/:interface/config", api.SetInterfaceConfig)
		network.POST("/interface/:interface/restart", api.RestartInterface)
		network.POST("/interface/:interface/:action", api.ControlInterface)

		// PPPoE
		network.GET("/interface/:interface/pppoe", api.GetPPPoEConfig)
		network.POST("/interface/:interface/pppoe", api.ConfigurePPPoE)

		// 代理
		network.GET("/proxy", api.GetProxyConfig)
		network.POST("/proxy", api.SetProxyConfig)

		// Wi-Fi
		network.GET("/wifi/scan", api.ScanWiFiNetworks)
		network.POST("/wifi/connect", api.ConnectToWiFi)
		network.POST("/wifi/disconnect", api.DisconnectWiFi)
		network.GET("/wifi/current", api.GetCurrentWiFiConnection)

		// DNS
		network.GET("/dns", api.GetDNSConfig)
		network.POST("/dns", api.SetDNSConfig)

		// IP
		network.PUT("/ip", api.UpdateIPConfig)
	}
}

// registerFirewall 防火墙规则与配置路由。
func registerFirewall(g *gin.RouterGroup) {
	firewall := g.Group("/security/firewall")
	requireAuth(firewall)
	{
		firewall.GET("/rules", api.GetFirewallRules)
		firewall.POST("/rules", api.CreateFirewallRule)
		firewall.PUT("/rules/:id", api.UpdateFirewallRule)
		firewall.DELETE("/rules/:id", api.DeleteFirewallRule)
		firewall.POST("/apply", api.ApplyFirewallRules)
		firewall.GET("/config", api.GetFirewallConfig)
		firewall.PUT("/config", api.SetFirewallConfig)
	}
}
