import React, {Component} from 'react';
import "./signup.css";
import {Link} from 'react-router-dom';
import axios from 'axios';



class SignUp extends Component {
    constructor(props){
        super(props);
        this.state = {
            fname : "",
            lname : "",
            email : "",
            password : ""
        }
        this.onChangeSignUp = this.onChangeSignUp.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
    }
    onChangeSignUp (e) {
        e.preventDefault();
        this.setState({[e.target.name] : e.target.value})
    }

     async onSubmit(e){
        e.preventDefault();
         const signUpData = {
            Firstname: this.state.fname,
            Lastname: this.state.lname,
            Email : this.state.email,
            Password : this.state.password
        }
        try{
            const connectionReqResponse = await axios.post('http://kong-elb-234657806.us-west-1.elb.amazonaws.com:80/user/users/signup', signUpData)
            if (connectionReqResponse.status === 201){
                alert("User has been succesfully created");
                this.props.history.push("/login");
            }
        } catch(err) {
            if (err.response.status === 409){
               alert("User with same email already exists")
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

            <div className="signupbox">
            <h1 className = "signupheading">CREATE ACCOUNT</h1>
            <center>Already having an account? <Link to="/login"> Sign in!</Link></center>
            <form onSubmit = {this.onSubmit}>
            <div className = "signUpDiv">
            <label for="FirstName" className="signUpLabel">
             First Name
            </label>
            <input type="text" name = "fname" onChange = {this.onChangeSignUp} className="signUpInput" required/>
            </div>

            <div className = "signUpDiv">
            <label for="LastName"  className="signUpLabel">
             Last Name
            </label>
            <input type="text" name = "lname" onChange = {this.onChangeSignUp} className="signUpInput" required/>
            </div>

            <div className = "signUpDiv">
            <label for="email" className="signUpLabel">
             Email Address
            </label>
            <input type="email" name = "email" onChange = {this.onChangeSignUp} className="signUpInput" required/>
            </div>

            <div className = "signUpDiv">
            <label for="Password" className="signUpLabel">
             Password
            </label>
            <input type="password" name = "password" onChange = {this.onChangeSignUp} className="signUpInput" required/>
            </div>

            <div className = "signUpDiv">
            <input type="submit" className="signUpButton" value="Create Account"/>
            </div>
            </form>
            </div>
            </div>
        )
    }
}
export default SignUp;
