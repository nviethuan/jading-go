// bids [
//         {
//                 "Price": "8.13800000",
//                 "Quantity": "456.02000000"
//         },
//         {
//                 "Price": "8.13700000",
//                 "Quantity": "446.99000000"
//         },
//         {
//                 "Price": "8.13600000",
//                 "Quantity": "631.61000000"
//         },
//         {
//                 "Price": "8.13500000",
//                 "Quantity": "893.60000000"
//         },
//         {
//                 "Price": "8.13400000",
//                 "Quantity": "1059.22000000"
//         },
//         {
//                 "Price": "8.13300000",
//                 "Quantity": "1251.94000000"
//         },
//         {
//                 "Price": "8.13200000",
//                 "Quantity": "849.55000000"
//         },
//         {
//                 "Price": "8.13100000",
//                 "Quantity": "1253.73000000"
//         },
//         {
//                 "Price": "8.13000000",
//                 "Quantity": "274.80000000"
//         },
//         {
//                 "Price": "8.12900000",
//                 "Quantity": "395.09000000"
//         },
//         {
//                 "Price": "8.12800000",
//                 "Quantity": "808.54000000"
//         },
//         {
//                 "Price": "8.12700000",
//                 "Quantity": "774.55000000"
//         },
//         {
//                 "Price": "8.12600000",
//                 "Quantity": "1420.84000000"
//         },
//         {
//                 "Price": "8.12500000",
//                 "Quantity": "1262.56000000"
//         },
//         {
//                 "Price": "8.12400000",
//                 "Quantity": "337.18000000"
//         },
//         {
//                 "Price": "8.12300000",
//                 "Quantity": "117.32000000"
//         },
//         {
//                 "Price": "8.12200000",
//                 "Quantity": "732.42000000"
//         },
//         {
//                 "Price": "8.12100000",
//                 "Quantity": "1048.68000000"
//         },
//         {
//                 "Price": "8.12000000",
//                 "Quantity": "582.51000000"
//         },
//         {
//                 "Price": "8.11900000",
//                 "Quantity": "42.66000000"
//         }
// ]
// asks [
//         {
//                 "Price": "8.13900000",
//                 "Quantity": "66.31000000"
//         },
//         {
//                 "Price": "8.14000000",
//                 "Quantity": "105.57000000"
//         },
//         {
//                 "Price": "8.14100000",
//                 "Quantity": "501.67000000"
//         },
//         {
//                 "Price": "8.14200000",
//                 "Quantity": "312.85000000"
//         },
//         {
//                 "Price": "8.14300000",
//                 "Quantity": "1172.08000000"
//         },
//         {
//                 "Price": "8.14400000",
//                 "Quantity": "1435.83000000"
//         },
//         {
//                 "Price": "8.14500000",
//                 "Quantity": "448.06000000"
//         },
//         {
//                 "Price": "8.14600000",
//                 "Quantity": "816.07000000"
//         },
//         {
//                 "Price": "8.14700000",
//                 "Quantity": "1166.99000000"
//         },
//         {
//                 "Price": "8.14800000",
//                 "Quantity": "649.13000000"
//         },
//         {
//                 "Price": "8.14900000",
//                 "Quantity": "1037.96000000"
//         },
//         {
//                 "Price": "8.15000000",
//                 "Quantity": "1824.69000000"
//         },
//         {
//                 "Price": "8.15100000",
//                 "Quantity": "306.05000000"
//         },
//         {
//                 "Price": "8.15200000",
//                 "Quantity": "276.53000000"
//         },
//         {
//                 "Price": "8.15300000",
//                 "Quantity": "783.05000000"
//         },
//         {
//                 "Price": "8.15400000",
//                 "Quantity": "2471.57000000"
//         },
//         {
//                 "Price": "8.15500000",
//                 "Quantity": "549.95000000"
//         },
//         {
//                 "Price": "8.15600000",
//                 "Quantity": "278.59000000"
//         },
//         {
//                 "Price": "8.15700000",
//                 "Quantity": "309.63000000"
//         },
//         {
//                 "Price": "8.15800000",
//                 "Quantity": "365.87000000"
//         }
// ]
package main

import (
	"fmt"

	binance "github.com/binance/binance-connector-go"
)

func wsDepthHandler() func(event *binance.WsPartialDepthEvent) {
	return func(event *binance.WsPartialDepthEvent) {
		fmt.Println("bids", binance.PrettyPrint(event.Bids))
		fmt.Println("asks", binance.PrettyPrint(event.Asks))
	}
}

func errHandler(err error) {
	fmt.Println(err)
}

func CheckStreamDepth() {
	websocketStreamClient := binance.NewWebsocketStreamClient(false)
	doneCh, _, err := websocketStreamClient.WsPartialDepthServe("UNIUSDT", "20", wsDepthHandler(), errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneCh
}
