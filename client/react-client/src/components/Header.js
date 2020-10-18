import React, { useState } from 'react';
import {
    Jumbotron, 
    Container, 
    Nav } from 'react-bootstrap';
import { Redirect } from 'react-router-dom';

export default function Header(props) {
    const [redirect, setRedirect] = useState(false);
    const logout = () => {
        localStorage.removeItem('username');
        localStorage.removeItem('token');
        localStorage.removeItem('expired');
        setRedirect(true);
    };

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
                    <Nav.Link onClick={logout} disabled={props.disabledLogout}>Logout</Nav.Link>
            </Nav.Item>
            </Nav>
        </div>
    );
}
