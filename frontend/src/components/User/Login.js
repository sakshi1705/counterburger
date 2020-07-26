import React, {Component} from 'react';
import "./signup.css";
import {Link} from 'react-router-dom';
import axios from 'axios';


class Login extends Component {
    constructor(props){
        super(props);
        this.state = {
            email : "",
            password : ""
        }
        this.onChangeLogin = this.onChangeLogin.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
    }
    onChangeLogin (e) {
        e.preventDefault();
        this.setState({[e.target.name] : e.target.value})
    }
     async onSubmit(e){
        e.preventDefault();
        console.log("in submit")
        const loginData = {
            Email : this.state.email,
            Password : this.state.password
        }
        try{
            const connectionReqResponse = await axios.post('http://kong-elb-234657806.us-west-1.elb.amazonaws.com:80/user/users/signin', loginData)
            if (connectionReqResponse.status === 201){
                alert("Login successful");
                let user = {
                    email:connectionReqResponse.data.email,
                    id: connectionReqResponse.data.id,
                    name: connectionReqResponse.data.firstname
                };
                localStorage.setItem('user', JSON.stringify(user));
                this.props.history.push("/location");
            }
            }
         catch(err) {
            if (err.response.status === 401){
                alert(err.response.data.Message)
         }
        }
     }
    render(){
        return(
            <div>
            <div className="signup"> 
            <div className="signupNav">
            <span className="signupSpan"> <a href="/" id="A_4"></a></span>
            </div>
            </div>


            <div className="loginbox">
            <h1 className = "signupheading">Login</h1>
            <center>Don't have an account?<Link to="/signup"> Create One!</Link> </center>
            <form onSubmit = {this.onSubmit}>
            <div className = "signUpDiv">
            <label for="email" className="signUpLabel">
             Email Address
            </label>
            <input type="email" name = "email" onChange = {this.onChangeLogin} class="signUpInput"/>
            </div>

            <div className = "signUpDiv">
            <label for="Password" className="signUpLabel">
             Password
            </label>
            <input type="password" name = "password" onChange = {this.onChangeLogin} className="signUpInput"/>
            </div>

            <div className = "signUpDiv">
            <input type="submit" className="signUpButton"  value="Login"/>
            </div>
            </form>
            </div>
            </div>
        )
    }


}
export default Login;
