import { Form, Input, message, notification } from "antd";
import { useNavigate } from "react-router-dom";
import {fetchWithAuth} from "../../utils/fetchWithAuth.js";

const CreateProject = ({ onSuccess }) => {
    const navigate = useNavigate();
    const [api, contextHolder] = notification.useNotification();

    const onFinish = async (values) => {
        try {
            const response = await fetchWithAuth("/api/v1/projects", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(values),
            }, navigate);

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
                name="createProject"
                onFinish={onFinish}
            >
                <Form.Item
                    label="Name"
                    name="name"
                    rules={[{ required: true }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="Description"
                    name="description"
                    rules={[{ required: true }]}
                >
                    <Input />
                </Form.Item>
            </Form>
        </>
    )
}

export default CreateProject;