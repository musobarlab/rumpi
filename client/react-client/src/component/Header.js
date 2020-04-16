import React, {Component} from 'react'
import {
    Jumbotron, 
    Container, 
    Nav } from 'react-bootstrap';

class Header extends Component {
  render() {
    return (
        <div>
            <Jumbotron fluid>
            <Container>
                <h1>Random Chat Demo</h1>
                <p>
                Random Chat Demo Using Golang and React.
                </p>
            </Container>
            </Jumbotron>
            <Nav className="justify-content-center" activeKey="/">
            <Nav.Item>
                <Nav.Link href="/">Home</Nav.Link>
            </Nav.Item>
            <Nav.Item>
                <Nav.Link href="/chat">Chat</Nav.Link>
            </Nav.Item>
            <Nav.Item>
                <Nav.Link href="/logout" disabled={this.props.disabledLogout}>Logout</Nav.Link>
            </Nav.Item>
            </Nav>
        </div>
    );
  }
}

export default Header;
