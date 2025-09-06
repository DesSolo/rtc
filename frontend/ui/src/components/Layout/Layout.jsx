import React, {useState} from "react";
import {Outlet, useNavigate} from 'react-router-dom';
import { Button, Layout, theme } from "antd";
import { MenuFoldOutlined, MenuUnfoldOutlined, LogoutOutlined } from "@ant-design/icons";
import NavBar from "../Nav/NavBar.jsx";

const { Header, Content } = Layout;

export const LayoutWithNav = () => {
    const [collapsed, setCollapsed] = useState(false);
    const [title, setTitle] = useState('');
    const navigate = useNavigate();
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    const handleLogout = () => {
        localStorage.removeItem('token')
        navigate('/login')
    }

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <NavBar collapsed={collapsed} />
            <Layout>
                <Header style={{
                    padding: 0,
                    background: colorBgContainer,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between'
                }}>
                    <div style={{ display: 'flex', alignItems: 'center' }}>
                        <Button
                            type="text"
                            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                            onClick={() => setCollapsed(!collapsed)}
                            style={{
                                fontSize: '16px',
                                width: 64,
                                height: 64,
                            }}
                        />
                        <span style={{ marginLeft: 16, fontSize: '18px', fontWeight: 'bold' }}>
                            {title}
                        </span>
                    </div>

                    <Button
                        type="text"
                        icon={<LogoutOutlined />}
                        onClick={handleLogout}
                        style={{
                            fontSize: '16px',
                            width: 64,
                            height: 64,
                        }}
                    />
                </Header>
                <Content
                    style={{
                        margin: '24px 16px',
                        padding: 24,
                        minHeight: 280,
                        background: colorBgContainer,
                        borderRadius: borderRadiusLG,
                    }}
                >
                    <Outlet context={{ setTitle }} />
                </Content>
            </Layout>
        </Layout>
    )
}

export default LayoutWithNav;