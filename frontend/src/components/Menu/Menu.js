import React, { Component } from 'react';
import counterburgersymbol from './counterburgersymbol.png';
import cbsymbol from './cbsymbol.jpg';
import burgerdetails from './burgerdetails.png';
import {Link} from 'react-router-dom';
import axios from 'axios';
import {Redirect} from 'react-router';
//import {Redirect} from 'react-router';
import './Menu.css';
import { Card } from 'antd';
//const { Meta } = Card;
var swal = require('sweetalert')
var hostname = 'http://kong-elb-234657806.us-west-1.elb.amazonaws.com:80/menu/menu' 
var hostnameOrder = 'http://kong-elb-234657806.us-west-1.elb.amazonaws.com:80/order/order' 
class Menu extends Component {
    
    constructor(props){
        super(props);
        this.state = {
            allmenu: [],  
            redirectVar : null
        }
        //this.handleLogout = this.handleLogout.bind(this);
        this.addToCart=this.addToCart.bind(this)
    }  

    addToCart(e1,e2,e3,e4,e5) {
        var user=JSON.parse(localStorage.getItem('user'));
        console.log(user.id);
        var headers = new Headers();
        //prevent page from refresh
        //e.preventDefault();
        const data = {
            UserId: user.id, //LOCAL STORAGE
            ItemId: e1,
            ItemName: e2,
            Price: e3, 
            Description: e4,
            ItemType : e5
        }
        console.log("dataa",data)

        //set the with credentials to true
        //axios.defaults.withCredentials = true;

        axios.post(hostnameOrder,data)
        .then(response => {
            console.log("Status Code : ",response.status);
            console.log("Data Sent ",response.data);
            if(response.status === 200){
                swal("Item Added To The Cart","","success") 
            }else{
                swal("Apologies, the item could not be added.","Please Try Again!","error")
            }
        
        });
    }

    componentDidMount()
    {
        axios.get(hostname)
                    .then((response) => {
                        console.log("Response data", response.data)
                        this.setState({
                            allmenu : response.data
                        })
                });
                console.log("Checking menu details", this.state.allmenu)
    }

    render() {
        let redirectVar = null;
        if(!localStorage.getItem("user")){
            redirectVar = <Redirect to= "/home"/>
        }
        let wholemenu = this.state.allmenu.map((wholemenu,j) => {
        return(
            <div className ="Menu">
                <Card
                hoverable
                style={{ width: 300 }}
                cover={<img src = {cbsymbol} height="320" width="550" alt=""></img>}
                className="MenuCards"
                >
                <div className = "ItemDescription">
                <br></br>
                <p><b>Item Type : </b>{wholemenu.ItemType}</p>
               <p><b>Item Name : </b>{wholemenu.ItemName}</p>
               <p><b>Description :</b> {wholemenu.Description}</p>
               <p><b>Price : </b>{wholemenu.Price}$</p>
               </div>
               <button onClick={()=>this.addToCart(wholemenu.ItemId, wholemenu.ItemName, wholemenu.Price, wholemenu.Description, wholemenu.ItemType)} className="btn btn-danger cartButton ">Add to Cart</button>
                </Card>
            </div>
        )
        })

        return (
            <div>
            {redirectVar}
            <div className="backgroundwallimage">
                <div className="counterburgersymbol">
                    <img src = {counterburgersymbol} height="100" width="200" alt=""></img>
                    <div className="NavbarLinks">
                    &nbsp; &nbsp; <Link to="/home" style={{"font-size": "20px", "font-weight" : "800" , "color":"black", "background-color": "white" }}>HOME</Link> 
                    <Link to="/menu" style={{"font-size": "20px", "font-weight" : "800" , marginLeft: "20px", "color":"black", "background-color": "white"  }}>MENU</Link>
                    <Link to="/burgerOrder" style={{"font-size": "20px", "font-weight" : "800" , marginLeft: "20px","color":"black", "background-color": "white"  }}>CART</Link>
                    </div>
                    <div className="container MenuOustide">
                    <div className="storedetails">
                        &nbsp;&nbsp; <b style={{ "font-size": "40px", "font-weight" : "800" , marginBottom: "0px" }}>THE COUNTER</b>
                        <br></br>
                        &nbsp;&nbsp; <Link to="/location">Change Location</Link>
                       <br></br>
                       <p>&nbsp;&nbsp; Phone: (408) 423-9200</p>
                       <p> &nbsp;&nbsp; Pickup Hours: Open today 11am-10pm </p>
                       <p>&nbsp;&nbsp; Accepted Cards: Mastercard, American Express, Discover</p>
                    
                    <div className="burgerdetails">
                    <img src = {burgerdetails} height="220" width="550" alt=""></img>
                    <p style={{"font-size": "15px", "font-weight" : "400",marginLeft: "5px", marginBottom:'5px'}}>Selections vary by location and may have limited availability.<br></br>
                        <a className="allergy" href="/nutrition">Nutritional information<br></br>
                        Allergen information</a>
                    </p>

                    </div>
                    </div>
                    <div className="container Menu2">
                        {wholemenu}
                        <br></br>
                    </div>
                </div>
                </div>
                </div>
            </div>
        );
    }
}

export default Menu;