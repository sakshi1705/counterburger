import React, { Component } from 'react';
import counterburgersymbol from './counterburgersymbol.png';
import {Link} from 'react-router-dom';
import axios from 'axios';
import './Payments.css';
import {Redirect} from 'react-router';

class Payments extends Component {

    constructor(props){
        super(props);
        this.state = {
            orderId : '',
            price: '',
            cardNumber: '',
            expirationDate: '',
            cvv: '',
            zipcode: '',
            country: ''
        }
        this.handleLogout = this.handleLogout.bind(this);
        this.onChange = this.onChange.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
    }
    async onSubmit(e){
        e.preventDefault();
        console.log("in onsubmit")
        const orderID = localStorage.getItem('orderId')
        const user = JSON.parse(localStorage.getItem('user'))
        const totalPrice = JSON.parse(localStorage.getItem('price'))
        const paymentsData = {
            OrderID : orderID,
            UserID : user.id,
            TotalPrice : totalPrice
        }
        try{
            const connectionReqResponse = await axios.post('http://kong-elb-234657806.us-west-1.elb.amazonaws.com:80/payment/payments', paymentsData)
            if (connectionReqResponse.status === 200){
                const orderDel = await axios.delete(`http://kong-elb-234657806.us-west-1.elb.amazonaws.com:80/order/order/${orderID}`)
                alert("Payment successful, Please LogOut!");
            }
            }
         catch(err) {
            window.alert(err.code)
            if (err.response.status === 404){
         }


    }
}
componentDidMount(){
    const orderID = localStorage.getItem('orderId')
    const totalPrice = localStorage.getItem('price')
    this.setState({
        orderId : orderID,
        price : totalPrice
    })
}


    handleLogout = () => {
        localStorage.removeItem('user');
        localStorage.removeItem('orderId');
        localStorage.removeItem('price');    
    }

    onChange(e){
        this.setState({[e.target.name]:e.target.value})
    }

    render() {
        let redirectVar = null;
        if(!localStorage.getItem("user")){
            redirectVar = <Redirect to= "/home"/>
        }
        return (
            <div>
            {redirectVar}
                <div className="counterburgersymbol">
                    <img src = {counterburgersymbol} height="100" width="200" alt=""></img>
                </div>
                <div className="logoutdiv">
                <li><Link to="/home" onClick = {this.handleLogout}><span class="glyphicon glyphicon-user logout"></span>Logout</Link></li>
                </div>
                <div className="container paymentcontainer">
                <br></br>
                    <b className="BillDetails">
                        Bill Details <br></br>
                    </b>
                    <b>
                        Order Id : {this.state.orderId} <br></br>
                        Amount to be Paid : {this.state.price} <br></br>
                    </b>
                    <br></br>
                    <p className="paymentMethod">
                        Payment Method - 
                    </p>
                    <p>
                        Enter Credit Card Details:
                    </p>
                    <form onSubmit = {this.onSubmit}>
                        Enter Card Number:
                        <div class="form-group">
                            <input onChange = {this.onChange} type="text" class="form-control" name="cardNumber" value={this.state.cardNumber} placeholder="Card Number" required/>
                        </div>
                        Enter Expiration Date:
                        <div class="form-group">
                            <input onChange = {this.onChange} type="date" class="form-control" name="expirationDate" value={this.state.expirationDate} placeholder="Expiration Date" required/>
                        </div>
                        Enter Security Code:
                        <div class="form-group">
                            <input onChange = {this.onChange} type="number" class="form-control" value={this.state.cvv} name="cvv" placeholder="CVV" required/>
                        </div>
                        Enter Zip/Postal Code:
                        <div class="form-group">
                            <input onChange = {this.onChange} type="number" class="form-control" value={this.state.zipcode} name="zipcode" placeholder="Zip Code" required/>
                        </div>
                        Select Country:
                        <div class="form-group">
                        <input onChange = {this.onChange} type="text" class="form-control" value={this.state.country} name="country" placeholder="Country" required/>
                        </div>
                        <div>
                            <input type="submit" class="btn btn-success btn-lg btn-block" value = "Place Order" />  
                        </div>
                    </form> 

                </div>
            </div>
        );
    }
}

export default Payments; 
