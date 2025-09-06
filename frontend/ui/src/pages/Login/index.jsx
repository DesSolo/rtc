import { Button, Checkbox, Form, Input, message, Card, Typography } from 'antd';
import {useEffect, useState} from "react";
import { useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

const { Title } = Typography;


const Login = () => {
    const [loading, setLoading] = useState(false);
    const [messageApi, contextHolder] = message.useMessage();
    const navigate = useNavigate();

    const onFinish = async (values) => {
        setLoading(true);
        try {
            const response = await fetch('/api/v1/login', {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username: values.username,
                    password: values.password
                }),
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem('token', data.data.token);

                const decoded = jwtDecode(data.data.token)
                localStorage.setItem('jwt', JSON.stringify(decoded))
                messageApi.success('Login successful!');
                navigate("/")
            } else if (response.status === 401) {
                messageApi.error('Username or password is incorrect');
            } else {
                messageApi.error('Login failed. Please try again.');
            }
        } catch (error) {
            messageApi.error('Network error. Please try again.');
        } finally {
            setLoading(false);
        }
    };

    const onFinishFailed = (errorInfo) => {
        console.log('Failed:', errorInfo);
    };

    return (
        <div style={{
            minHeight: '100vh',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)'
        }}>
            {contextHolder}
            <Card
                style={{
                    width: 400,
                    boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
                    borderRadius: 8
                }}
                bodyStyle={{ padding: 32 }}
            >
                <div style={{ textAlign: 'center', marginBottom: 24 }}>
                    <Title level={2} style={{ margin: 0, color: '#1890ff' }}>Sign In</Title>
                    <p style={{ marginTop: 8, color: '#8c8c8c' }}>Enter your credentials to continue</p>
                </div>

                <Form
                    name="basic"
                    layout="vertical"
                    initialValues={{ remember: true }}
                    onFinish={onFinish}
                    onFinishFailed={onFinishFailed}
                    autoComplete="off"
                    disabled={loading}
                >
                    <Form.Item
                        label="Username"
                        name="username"
                        rules={[{ required: true, message: 'Please input your username!' }]}
                    >
                        <Input size="large" placeholder="Enter your username" />
                    </Form.Item>

                    <Form.Item
                        label="Password"
                        name="password"
                        rules={[{ required: true, message: 'Please input your password!' }]}
                    >
                        <Input.Password size="large" placeholder="Enter your password" />
                    </Form.Item>

                    <Form.Item name="remember" valuePropName="checked">
                        <Checkbox>Remember me</Checkbox>
                    </Form.Item>

                    <Form.Item>
                        <Button
                            type="primary"
                            htmlType="submit"
                            loading={loading}
                            size="large"
                            block
                        >
                            Sign In
                        </Button>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    );
};

export default Login;