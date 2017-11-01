//
// Date: 10/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
  "github.com/tidwall/gjson"   
)

type Balance struct {
  AccountNumber string
  AccountValue float64
  TotalCash float64 
  OptionBuyingPower float64
  StockBuyingPower float64
}

//
// Get Balances
//
func GetBalances() ([]Balance, error) {
  
  var balances []Balance
  
  // Get the JSON
  jsonRt, err := SendGetRequest("/user/balances")

  if err != nil {
    return balances, err  
  } 
  
  // Make sure we have at least one account (this should never happen)
  vo := gjson.Get(jsonRt, "accounts.account")
  
  if ! vo.Exists() {
    return balances, nil  
  }  
  
  // Do we have only one account?
  vo = gjson.Get(jsonRt, "accounts.account.balances")
  
  // Only one account
  if vo.Exists() {
      
    // Add to balances array
    balances = append(balances, addJsonToBalance(vo.String()))
           
  } else // More than one accounts
  {

    vo := gjson.Get(jsonRt, "accounts.account")
  
    // Loop through the different accounts.
    vo.ForEach(func(key, value gjson.Result) bool {
      
      // Add to balances array
      balances = append(balances, addJsonToBalance(gjson.Get(value.String(), "balances").String()))
      
      // keep iterating
      return true  
          
    }) 
        
  } 
  
  // Return happy
  return balances, nil
  
}

//
// Take some JSON and return a Balance object.
//
func addJsonToBalance(accountJson string) Balance {
  
  var optionBuyingPower float64
  var stockBuyingPower float64
  
  // Get option buying power
  vo3 := gjson.Get(accountJson, "margin.option_buying_power")
              
  if vo3.Exists() {
    
    optionBuyingPower = vo3.Float()
  
  } else {
    
    optionBuyingPower = gjson.Get(accountJson, "cash.cash_available").Float() 
           
  }
  
  // Get stock buying power
  vo3 = gjson.Get(accountJson, "margin.stock_buying_power")
              
  if vo3.Exists() {
    
    stockBuyingPower = vo3.Float()
  
  } else {
    
    stockBuyingPower = gjson.Get(accountJson, "cash.cash_available").Float() 
           
  }
              
  // Return balance
  return Balance{
    AccountNumber: gjson.Get(accountJson, "account_number").String(),
    AccountValue: gjson.Get(accountJson, "total_equity").Float(),
    TotalCash: gjson.Get(accountJson, "total_cash").Float(),
    OptionBuyingPower: optionBuyingPower,
    StockBuyingPower: stockBuyingPower,
  }
  
}