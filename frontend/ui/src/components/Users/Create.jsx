import {Button, Checkbox, Form, Input, message, notification, Space} from "antd";
import {fetchWithAuth} from "../../utils/fetchWithAuth.js";
import {PlusOutlined, MinusCircleOutlined} from "@ant-design/icons";

const CreateUser = ({ onSuccess }) => {
    const [api, contextHolder] = notification.useNotification();

    const onFinish = async (values) => {
        try {
            const response = await fetchWithAuth("/api/v1/users", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(values),
            });

            if (!response.ok) {
                const errorData = await response.json();
                api.error({
                    message:  "Error",
                    description: errorData.error.toString() || "Some error"
                })
                return
            }

            api.success({
                message: "Project created successfully"
            });

            if (onSuccess) {
                onSuccess();
            }
        } catch (error) {
            api.error({
                message: error.message
            });
        }
    };

    return (
        <>
            {contextHolder}
            <Form
                name="createUser"
                onFinish={onFinish}
            >
                <Form.Item
                    label="Username"
                    name="username"
                    rules={[{ required: true }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="Password"
                    name="password"
                    rules={[{ required: true }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item name="is_enabled" valuePropName="checked" label={null}>
                    <Checkbox>Enabled</Checkbox>
                </Form.Item>
                <Form.List name="roles">
                    {(fields, { add, remove }) => (
                        <>
                            {fields.map(({ key, name, ...restField }) => (
                                <Space key={key} align="baseline">
                                    <Form.Item
                                        {...restField}
                                        name={name}
                                    >
                                        <Input placeholder="Role name" />
                                    </Form.Item>
                                    <MinusCircleOutlined onClick={() => remove(name)} />
                                </Space>
                            ))}
                            <Form.Item>
                                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                                    Add role
                                </Button>
                            </Form.Item>
                        </>
                    )}
                </Form.List>

            </Form>
        </>
    )
}

export default CreateUser;