import React, {Component} from 'react'
import {
    Jumbotron, 
    Container, 
    Nav } from 'react-bootstrap';
import {Redirect} from 'react-router-dom';

class Header extends Component {

    constructor(props) {
        super(props);

        this.state = {
            redirect: false
        }

        this._logout = this._logout.bind(this);
    }

    _logout() {
        localStorage.removeItem('username');
        localStorage.removeItem('token');
        localStorage.removeItem('expired');
        this.setState({redirect: true});
    }

    render() {
        const {redirect} = this.state;

        if (redirect) {
            return <Redirect to='/'/>;
        }
        return (
            <div>
                <Jumbotron fluid>
                    <Container>
                        <h1>Rumpi</h1>
                        <p>
                        Opensource chat application scaffold for building chat application with Go and Reactjs
                        </p>
                    </Container>
                </Jumbotron>
                <Nav className="justify-content-center" activeKey="/">
                    <Nav.Item>
                        <Nav.Link href="/">Home</Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link href="/register">Register</Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link href="/chat">Chat</Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link onClick={this._logout} disabled={this.props.disabledLogout}>Logout</Nav.Link>
                </Nav.Item>
                </Nav>
            </div>
        );
    }
}

export default Header;
