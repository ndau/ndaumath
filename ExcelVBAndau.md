'
'   ndau Price Calculation Functions 1.0
'
'   Consider these the master definitions of these functions and the methods by which these values are calculated.
'   Use these functions as pseudocode for implementation in other situations; the code is not compact and variable
'   names are quite verbose to make it more helpful for that purpose.
'
'   All functions and calculations are only defined for Phase 1.
'
'   BLOCKS ARE NUMBERED BEGINNING AT 0. THE FIRST 1,000 NDAU ARE IN BLOCK 0.
'
'   System Constants - don't assume anything is immutable, starting with the price of the very first ndau issued
'
Const NdauStartPrice = 1
Const BlockSize = 1000
Const Phase1Blocks = 10000
Const Phase1Doublings = 14
Const Phase1BlockPriceRatio = 2 ^ (Phase1Doublings / Phase1Blocks)
'
' To calculate the price at which the last ndau were sold, we need to subtract just a tiny bit from the total number sold (to avoid block boundary conditions)
'
Const Epsilon = 1E-09
'
'   Excel has a built-in MIN() function, but VBA does not.
'
Function Min(a, b)
    If a > b Then Min = b Else Min = a
End Function

'
'   ndau_block_number
'
'   Find the block number in which the specified ndau falls (isolates BlockSize)
'
Function ndau_block_number(ndau)
    ndau_block_number = Int(ndau / BlockSize)
End Function
'
'   ndau_price
'
'   Find the ndau issuance price at any sales level. By default it's the price at which the next ndau will be sold, but it can be the
'   price at which the last ndau were sold. The latter is more suitable for market cap calculation.
'
Function ndau_price(ndau_issued, Optional next_or_last = "next")
    If next_or_last = "next" Then
        ndau_price = NdauStartPrice * (Phase1BlockPriceRatio ^ ndau_block_number(ndau_issued))
    Else
         ndau_price = NdauStartPrice * (Phase1BlockPriceRatio ^ ndau_block_number(ndau_issued - Epsilon))
        End If
End Function
'
'   ndau_block_price
'
'   Return the price of ndau in the Nth BlockSize-sized block.
'
Function ndau_block_price(block_number)
    ndau_block_price = Phase1BlockPriceRatio ^ block_number
End Function
'   ndau_cost_blocks
'
'   The cost of a specific number of whole blocks of ndau beginning at a given block.
'
Function ndau_cost_blocks(starting_block, number_of_blocks)
    ndau_cost_blocks = NdauStartPrice * BlockSize * (ndau_block_price(starting_block + number_of_blocks) - ndau_block_price(starting_block)) / (Phase1BlockPriceRatio - 1)
End Function
'
'   ndau_cost
'
'   Total dollar cost of ndau after N have been issued from the endowment
'
Function ndau_cost(ndau_issued, count)
    first_block = ndau_block_number(ndau_issued)
    ndau_available_in_first_block = BlockSize - (ndau_issued - (first_block * BlockSize))
    ndau_available_in_first_block = BlockSize - ndau_issued Mod BlockSize
    
    ndau_sold_from_first_block = Min(ndau_available_in_first_block, count)
    cost = ndau_sold_from_first_block * ndau_block_price(first_block)
    
    ndau_remaining = count - ndau_sold_from_first_block
    ndau_remaining_whole_blocks = Int(ndau_remaining / BlockSize)
    cost = cost + ndau_cost_blocks(first_block + 1, ndau_remaining_whole_blocks)
    
    ndau_in_last_block = ndau_remaining - (ndau_remaining_whole_blocks * BlockSize)
    cost = cost + (ndau_in_last_block * ndau_block_price(first_block + ndau_remaining_whole_blocks + 1))
    
    ndau_cost = cost
End Function
'
'   ndau_count
'
'   Total number of ndau sold for a given dollar cost after N have been issued from the endowment
'
Function ndau_count(ndau_issued, cost)
    first_block = ndau_block_number(ndau_issued)
    ndau_in_first_block_price = ndau_block_price(first_block)
    
    ndau_available_in_first_block = BlockSize - (ndau_issued Mod BlockSize)
    ndau_in_first_block = Min(ndau_available_in_first_block, cost / ndau_in_first_block_price)
    ndau_in_first_block_cost = ndau_in_first_block * ndau_in_first_block_price
    money_remaining = cost - ndau_in_first_block_cost
    
    next_block = first_block + 1
    b1 = (money_remaining / BlockSize) * (Phase1BlockPriceRatio - 1)
    b2 = Phase1BlockPriceRatio ^ next_block
    blocks_to_be_purchased = Int(Log(b1 + b2) / Log(Phase1BlockPriceRatio)) - next_block
    ndau_in_whole_blocks = blocks_to_be_purchased * BlockSize
    money_remaining = money_remaining - ndau_cost(ndau_issued + ndau_in_first_block, ndau_in_whole_blocks)
   
    ndau_in_last_block = money_remaining / ndau_block_price(first_block + blocks_to_be_purchased + 1)
   
    ndau_count = ndau_in_first_block + ndau_in_whole_blocks + ndau_in_last_block
End Function
'
'   ndau_endowment_proceeds
'
'   The total proceeds delivered to the endowment from the sale of the specified number of ndau.
'
Function ndau_endowment_proceeds(ndau_issued)
    whole_block_proceeds = ndau_cost_blocks(0, ndau_block_number(ndau_issued))
    last_block_proceeds = (ndau_issued - (ndau_block_number(ndau_issued) * BlockSize)) * ndau_price(ndau_issued)
    ndau_endowment_proceeds = whole_block_proceeds + last_block_proceeds
End Function
'
'   ndau_market_cap
'
'   Total market cap of all ndau after a specified number have been issued from the endowment.
'
Function ndau_market_cap(ndau_issued)
    ndau_market_cap = ndau_issued * ndau_price(ndau_issued, "last")
End Function
'
'   ndau_sold_at_target_price
'
'   Given a current Target Price, return the number of ndau sold to reach that price.
'
Function ndau_sold_at_target_price(target_price)
    ndau_sold_at_target_price = Int(Log(target_price) / Log(Phase1BlockPriceRatio)) * 1000
End Function
