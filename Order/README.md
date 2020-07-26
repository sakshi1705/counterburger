## Order Module APIs

## APIs to manage Orders

1. Ping the API endpoint

    **Request**
    ```
    GET /orders/ping
    Content-Type: application/json
    ```
    **Response**
    ```
    {
    "vCloud9.0 Burger Order API version 1.0 alive!"
    }
    ```
    
    ---
    
2. Get all the orders:

    **Request**
    
    ```
    GET /order/{orderId}
    Content-Type: application/json
    ```

    **Response Body**

    (Status code: 200)

    |Parameter	|Type	|Description  |
    |----|----|----|
    |orderId |String | Order ID of the order placed |
    |userId | String | Id of user who has placed the order |
    |orderStatus | String  | Status of the order(Active/Placed)|
    |items |Struct | itemId,itemName,itemType,price,description |
    |totalAmount	| float32 | Total amount of the order placed by user|
   

    (Status code: 404)

    |Parameter	|Type |	Description|
    |-----|-----|------|
    |messsage	|String|Error string(No order found)|
    
    
3. Get order by the provided Order ID:

    **Request**
    
    ```
    GET /order/{orderId}
    Content-Type: application/json
    ```
    |Parameter	|Type |	Description|
    |-----|-----|------|
    |orderId	|String|Order Id of the placed order|

    **Response Body**

    (Status code: 200)

    |Parameter	|Type	|Description  |
    |----|----|----|
    |orderId |String | Order ID of the order placed |
    |userId | String | Id of user who has placed the order |
    |orderStatus | String  | Status of the order(Active/Placed)|
    |items |Struct | itemId,itemName,itemType,price,description |
    |totalAmount	| float32 | Total amount of the order placed by user|
   

    (Status code: 404)

    |Parameter	|Type |	Description|
    |-----|-----|------|
    |messsage	|String|Error string(Order with provided orderId not found)|
    
    ---
    
    
4.  Get order by the provided User ID:

    **Request**
    
    ```
    GET /orders/{userId}
    Content-Type: application/json
    ```
    |Parameter	|Type |	Description|
    |-----|-----|------|
    |userId	|String| Id of the user who has placed the order |

    **Response Body**

    (Status code: 200)

    |Parameter	|Type	|Description  |
    |----|----|----|
    |orderId |String | Order ID of the order placed |
    |userId | String | Id of user who has placed the order |
    |orderStatus | String  | Status of the order(Active/Placed)|
    |items |Struct | itemId,itemName,itemType,price,description |
    |totalAmount	| float32 | Total amount of the order placed by user|
   

    (Status code: 404)

    |Parameter	|Type |	Description|
    |-----|-----|------|
    |messsage	|String|Error string(Order with provided userId not found)|
    
    ---
    
    
5. Post a burger order done by a user:

   **Request**
    
    ```
    POST /order
    Content-Type: application/json
    ```
    
    |Parameter	|Type	|Description  |
    |----|----|----|
    |OrderId |String | Order ID of the order placed |
    |UserId | String | Id of user who has placed the order |
    |ItemType | String  | Type of order in burger order)|
    |ItemName | Struct | Name of the item |
    |Price	| float32 | Price of the item added to cart|
    |Description	| String | Description of the item|
    
     **Response Body**

    (Status code: 200)
    
    |Parameter	|Type	|Description  |
    |----|----|----|
    |orderId |String | Order ID of the order placed |
    |userId | String | Id of user who has placed the order |
    |orderStatus | String  | Status of the order:Active |
    |items |Struct | itemId,itemName,itemType,price,description |
    |totalAmount	| float32 | Total amount of the order placed by user|
    
    (Status code: 404)

    |Parameter	|Type |	Description|
    |-----|-----|------|
    |messsage	|String| Error string|
    
    ---
    


6. Update the order status after payment is successfully done:

    **Request**
    
    ```
    PUT /order/{orderId}
    Content-Type: application/json
    ```
    |Parameter	|Type |	Description|
    |-----|-----|------|
    | orderId | String | Id of the order placed by the user|
    
     **Response Body**
    
    (Status code: 200)
    
    |Parameter	|Type	|Description  |
    |----|----|----|
    |orderId |String | Order ID of the order placed |
    |userId | String | Id of user who has placed the order |
    |orderStatus | String  | Status of the order:Placed |
    |items |Struct | itemId,itemName,itemType,price,description |
    |totalAmount	| float32 | Total amount of the order|
    
    (Status code: 404)

    |Parameter	|Type |	Description|
    |-----|-----|------|
    |messsage	|String| Error string(Sorry!OrderId not found!) |
    
    ---
    


7. Delete an item from the order:

    **Request**
    
    ```
    DELETE /order/item/{orderId}
    Content-Type: application/json
    ```
    |Parameter	|Type |	Description|
    |-----|-----|------|
    | orderId | String | Id of the order placed by the user|
    
    #### Request Body:
    
    |Parameter	|Type	|Description  |
    |----|----|----|
    | itemId | String | Item Id of the item to be deleted |
    
    **Response Body**
    
    (Status code: 200)
    
    |Parameter	|Type	|Description  |
    |----|----|----|
    |orderId |String | Order ID of the order placed |
    |userId | String | Id of user who has placed the order |
    |orderStatus | String  | Status of the order:Active |
    |items |Struct | itemId,itemName,itemType,price,description |
    |totalAmount	| float32 | Total amount after the pice deduction of item deleted|
    
    (Status code: 404)

    |Parameter	|Type |	Description|
    |-----|-----|------|
    |messsage	|String| Error string(Sorry!OrderId not found!/ ItemId not found) |
    
    ---
   


8. Delete the entire order:

    **Request**
    
    ```
    DELETE /orders/{orderId}
    Content-Type: application/json
    ```
    |Parameter	|Type |	Description|
    |-----|-----|------|
    | orderId |String| Id of the order placed by the user|

    **Response Body**
    
    
    (Status code: 200)

    |Parameter	|Type |	Description|
    |-----|-----|------|
    |messsage	|String| Success message(Order has been deleted) |
    
    
    (Status code: 404)

    |Parameter	|Type |	Description|
    |-----|-----|------|
    |messsage	|String| Error string(Sorry!Order not found! |
    
    
