import { useState } from "react";
import { useNavigate } from 'react-router-dom';
import EnvironmentsSelector from "../Environments/Selector.jsx";
import { Table, Typography } from "antd";

const {Link} = Typography


const ReleasesList = ({project}) => {
    const navigate = useNavigate();
    const [env, setEnv] = useState();
    const [releases, setReleases] = useState([]);

    const columns = [
        {
            title: "Name",
            dataIndex: "name",
            key: "name",
        },
        {
            title: "Created",
            dataIndex: "created_at",
            key: "created",
            render: (text) => (
                <>
                    {new Date(text).toLocaleString()}
                </>
            )
        },
        {
            title: "Actions",
            key: "actions",
            render: (release) => (
                <>
                    <Link onClick={() => handleClick(release.name)}>conf</Link>
                </>
            )
        }
    ]

    const fetchRelease = (env) => {
        fetch(`/api/v1/projects/${project}/envs/${env}/releases`)
            .then((response) => {
                if (!response.ok) throw new Error("err");
                return response.json();
            })
            .then((data) => {
                setReleases(data.data.releases)
            })
    }

    const handleEnvChanged = (env) => {
        fetchRelease(env)
        setEnv(env)
    }

    const handleClick = (release) => {
        navigate(`/projects/${project}/envs/${env}/releases/${release}/configs`)
    }

    return (
        <>
            <EnvironmentsSelector project={project} onEnvChange={(selectedEnv) => (handleEnvChanged(selectedEnv)) }/>
            <Table
                columns={columns}
                dataSource={releases}
                rowKey="name"
                pagination={false}
            />
        </>
    )
}

export default ReleasesList;