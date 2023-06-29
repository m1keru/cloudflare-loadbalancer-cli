package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type cfLoadBalancerPool struct {
	Api       *cloudflare.API
	AccountID string
	ZoneID    string
	ctx       context.Context
}

func (pool cfLoadBalancerPool) updateCfLoadBalancerPool(poolName string, originName string, originState bool) {
	pools, err := pool.Api.ListLoadBalancerPools(pool.ctx, cloudflare.AccountIdentifier(pool.AccountID), cloudflare.ListLoadBalancerPoolParams{})
	for _, lbp := range pools {
		if lbp.Name == poolName {
			fmt.Printf("lbp.Name: %v\n", lbp.Name)
			var lbo_new []cloudflare.LoadBalancerOrigin
			for _, lbo := range lbp.Origins {
				fmt.Printf("	lbo: name: %s enabled: %t \n", lbo.Name, lbo.Enabled)
				if lbo.Name == originName && lbo.Enabled != originState {
					lbo.Enabled = originState
					fmt.Printf("		[change]: change state of origin %s to enabled:%t in:%s\n", lbo.Name, originState, lbp.Name)
				}
				lbo_new = append(lbo_new, lbo)
			}
			lbp.Origins = lbo_new
			_, err = pool.Api.UpdateLoadBalancerPool(pool.ctx, cloudflare.AccountIdentifier(pool.AccountID), cloudflare.UpdateLoadBalancerPoolParams{LoadBalancer: lbp})
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

		}
	}
}

func (pool cfLoadBalancerPool) listCfLoadBalancerPools() {
	pools, err := pool.Api.ListLoadBalancerPools(pool.ctx, cloudflare.AccountIdentifier(pool.AccountID), cloudflare.ListLoadBalancerPoolParams{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	//rc := cloudflare.ZoneIdentifier(zoneID)
	for _, lbp := range pools {
		fmt.Printf("lbp.Name: %v\n", lbp.Name)
		for _, lbo := range lbp.Origins {
			fmt.Printf("  lbo: name: %s enabled: %t \n", lbo.Name, lbo.Enabled)
		}
	}
}

func main() {

	list := flag.Bool("list", false, "list loadbalancers")
	update := flag.Bool("update", false, "update particular lb")
	poolName := flag.String("pool", "", "load balancer pool name")
	originName := flag.String("origin", "", "load balancer origin name")
	originState := flag.String("state", "", "load balancer origin state [enable|disable]")

	flag.Parse()

	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	accountID := "680bdb4c4755e7d9c0839b8377d9aa3e"
	zoneID, err := api.ZoneIDByName("easybrain.com")

	cli := cfLoadBalancerPool{
		Api:       api,
		AccountID: accountID,
		ZoneID:    zoneID,
		ctx:       ctx,
	}

	if *list {
		cli.listCfLoadBalancerPools()
		return
	}
	if *update {
		if *poolName == "" || *originName == "" || *originState == "" {
			fmt.Println("Error: define all params")
			flag.Usage()
			return
		}
		var state bool
		if *originState == "enable" {
			state = true
		} else if *originState == "disable" {
			state = false
		} else {
			fmt.Println("Error: wrong state. variants: [enable|disable]")
			return
		}
		cli.updateCfLoadBalancerPool(*poolName, *originName, state)
		return
	}

	flag.Usage()

}
