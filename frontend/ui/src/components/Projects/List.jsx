import {Button, Flex, Table} from "antd"
import { useEffect, useState } from "react";

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
            <span>{new Date(text).toLocaleString()}</span>
        )
    }
]

const ProjectsList = () => {
    const [data, setData] = useState([]);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 10,
        hasNextPage: false
    });

    const fetchData = async (page = 1, pageSize = 10) => {
        const response = await fetch (
            `/api/v1/projects?limit=${pageSize}&offset=${(page-1) * pageSize}`
        )

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

    useEffect(() => {
            fetchData()
    }, [])

    return (
        <>
            <Table
                columns={columns}
                dataSource={data}
                pagination={false}
            />
            <Flex align="center" justify="flex-end">
                <Button type="link" onClick={() => (stepPage(-1))} disabled={pagination.current === 1}>&#60;</Button>
                {pagination.current}
                <Button type="link" onClick={() => (stepPage(+1))} disabled={!pagination.hasNextPage}>&#x3e;</Button>
            </Flex>
        </>
    )
};

export default ProjectsList;