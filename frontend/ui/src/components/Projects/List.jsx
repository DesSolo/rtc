import {Button, Flex, Table, Modal, Input, Typography} from "antd"
import { useEffect, useState } from "react";
import { PlusOutlined } from '@ant-design/icons';
import CreateProject from "./Create.jsx";
import {useNavigate} from "react-router-dom";

const {Link} = Typography

const ProjectsList = () => {
    const navigate = useNavigate();
    const [data, setData] = useState([]);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 10,
        hasNextPage: false
    });
    const [isModalOpen, setIsModalOpen] = useState(false);

    const columns = [
        {
            title: "Name",
            dataIndex: "name",
            key: "name"
        },
        {
            title: "Description",
            dataIndex: "description",
            key: "description"
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
            render: (project) => (
                <>
                    <Link onClick={() => navigateToReleases(project.name)}>releases</Link>
                </>
            )
        }
    ]

    const navigateToReleases = (project) => {
        navigate(`/projects/${project}/releases`)
    }

    const fetchData = async (page = 1, pageSize = 10, q) => {
        let uri = "/api/v1/projects?"
        if (q) {
            uri += `&q=${q}`
        }
        uri += `&limit=${pageSize}&offset=${(page-1) * pageSize}`

        const response = await fetch (uri)

        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }

        const result = await response.json();
        const projects = result['data']['projects']

        setData(projects);

        setPagination({
            ...pagination,
            current: page,
            pageSize: pageSize,
            hasNextPage: projects.length === pageSize
        });
    }

    const stepPage = async (index) => {
        fetchData(pagination.current+index, pagination.pageSize)
    }

    const search = async (q) => {
        if (q.length === 0) {
            fetchData()
            return
        }

        if (q.length < 3) {
            return
        }

        fetchData(1, 10, q)
    }

    useEffect(() => {
            fetchData()
    }, [])

    return (
        <>
            <Flex justify="space-between" style={{ marginBottom: 16 }}>
                <Input placeholder="Search" onChange={(e) => search(e.target.value)} />
                <Button type="primary" icon={<PlusOutlined />} onClick={() => (setIsModalOpen(true))}>New</Button>
            </Flex>
            <Table
                columns={columns}
                dataSource={data}
                rowKey="name"
                pagination={false}
            />
            <Flex align="center" justify="flex-end" style={{ marginTop: 15 }}>
                <Button type="link" onClick={() => (stepPage(-1))} disabled={pagination.current === 1}>&#60;</Button>
                {pagination.current}
                <Button type="link" onClick={() => (stepPage(+1))} disabled={!pagination.hasNextPage}>&#x3e;</Button>
            </Flex>
            <Modal
                title="Create new project"
                open={isModalOpen}
                onCancel={()=>setIsModalOpen(false)}
                footer={[
                    <Button form="createProject" type="primary" key="submit" htmlType="submit">Create</Button>
                ]}
            >
                <div style={{marginTop: 20}}>
                    <CreateProject onSuccess={fetchData}/>
                </div>
            </Modal>
        </>
    )
};

export default ProjectsList;