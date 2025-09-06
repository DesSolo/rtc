import {useCallback, useEffect, useState} from "react";
import {useNavigate, useOutletContext} from "react-router-dom";
import EnvironmentsSelector from "../Environments/Selector.jsx";
import { Button, Table, Flex, Popconfirm, notification } from "antd";
import { DeleteOutlined } from "@ant-design/icons";
import {fetchWithAuth} from "../../utils/fetchWithAuth.js";

const ReleasesList = ({ project }) => {
    const { setTitle } = useOutletContext();
    const navigate = useNavigate();
    const [env, setEnv] = useState();
    const [releases, setReleases] = useState([]);
    const [api, contextHolder] = notification.useNotification();

    useEffect(() => {
        setTitle(`Releases: ${project}`)
    }, []);

    const fetchReleases = useCallback(async (environment) => {
        try {
            const resp = await fetchWithAuth(`/api/v1/projects/${project}/envs/${environment}/releases`);
            if (!resp.ok) throw new Error(`status ${resp.status}`);
            const data = await resp.json();
            setReleases(data.data?.releases || []);
        } catch (err) {
            api.error({ message: "Failed to load releases", description: String(err) });
        }
    }, [project, api]);

    const deleteRelease = useCallback(async (name) => {
        try {
            const resp = await fetchWithAuth(`/api/v1/projects/${project}/envs/${env}/releases/${name}`, {
                method: "DELETE",
            });
            if (!resp.ok) throw new Error(`status ${resp.status}`);
            api.success({ message: `Release ${name} deleted` });
            fetchReleases(env);
        } catch (err) {
            api.error({ message: "Failed to delete release", description: String(err) });
        }
    }, [env, project, fetchReleases, api]);

    const handleEnvChanged = (environment) => {
        setEnv(environment);
        fetchReleases(environment);
    };

    const handleClick = (release) => {
        navigate(`/projects/${project}/envs/${env}/releases/${release}/configs`);
    };

    const columns = [
        {
            title: "Name",
            dataIndex: "name",
            key: "name",
            render: (text) => (
                <a onClick={() => handleClick(text)}>{text}</a>
            ),
        },
        {
            title: "Created",
            dataIndex: "created_at",
            key: "created",
            render: (text) => new Date(text).toLocaleString(),
        },
        {
            title: "Actions",
            key: "actions",
            render: (release) => (
                <Flex gap="middle">
                    <Popconfirm
                        title={`Delete release ${release.name}?`}
                        onConfirm={() => deleteRelease(release.name)}
                        okText="Yes"
                        cancelText="No"
                    >
                        <Button danger type="primary" icon={<DeleteOutlined />} />
                    </Popconfirm>
                </Flex>
            ),
        },
    ];

    return (
        <>
            {contextHolder}
            <div style={{ display: "flex", gap: 12, marginBottom: 12 }}>
                <EnvironmentsSelector project={project} onEnvChange={handleEnvChanged} />
            </div>
            <Table
                columns={columns}
                dataSource={releases}
                rowKey="name"
                pagination={false}
            />
        </>
    );
};

export default ReleasesList;
