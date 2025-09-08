import { Button, Flex, Table, Modal, Input } from "antd";
import { useEffect, useState, useCallback, useRef } from "react";
import { PlusOutlined } from "@ant-design/icons";
import CreateProject from "./Create.jsx";
import { useNavigate, useOutletContext, useSearchParams } from "react-router-dom";
import {fetchWithAuth} from "../../utils/fetchWithAuth.js";

const ProjectsList = () => {
    const { setTitle } = useOutletContext();
    const navigate = useNavigate();
    const [searchParams, setSearchParams] = useSearchParams();

    // Get initial parameters from URL
    const paramPage = parseInt(searchParams.get("page") || "1", 10);
    const paramLimit = parseInt(searchParams.get("limit") || "10", 10);
    const paramQ = searchParams.get("q") || "";

    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [pagination, setPagination] = useState({
        current: paramPage,
        pageSize: paramLimit,
        total: 0,
        showSizeChanger: true,
        pageSizeOptions: ['10','20','50','100'],
    });
    const [searchValue, setSearchValue] = useState(paramQ);
    const [isModalOpen, setIsModalOpen] = useState(false);

    // Track last applied search value to avoid unnecessary URL updates
    const lastAppliedQRef = useRef(paramQ);
    const isFirstMount = useRef(true);

    const fetchData = useCallback(async (page = 1, pageSize = 10, q = "") => {
        setLoading(true);
        try {
            const uri = `/api/v1/projects?limit=${pageSize}&offset=${(page - 1) * pageSize}${q ? `&q=${encodeURIComponent(q)}` : ""}`;
            const res = await fetchWithAuth(uri, {}, navigate);
            if (!res.ok) throw new Error(res.statusText);
            const json = await res.json();
            setData(json.data.projects || []);
            setPagination(prev => ({
                ...prev,
                current: page,
                pageSize,
                total: json.data.total || 0,
            }));
        } catch (e) {
            console.error(e);
        } finally {
            setLoading(false);
        }
    }, []);

    // Sync with URL parameters
    useEffect(() => {
        setPagination(p => ({ ...p, current: paramPage, pageSize: paramLimit }));
        setSearchValue(paramQ);
        lastAppliedQRef.current = paramQ;
        fetchData(paramPage, paramLimit, paramQ);
        isFirstMount.current = false;
    }, [searchParams.toString()]);

    const handleTableChange = (newPagination) => {
        const nextPage = newPagination.current;
        const nextLimit = newPagination.pageSize;
        setSearchParams({ page: String(nextPage), limit: String(nextLimit), q: lastAppliedQRef.current || "" }, { replace: true });
        fetchData(nextPage, nextLimit, lastAppliedQRef.current || "");
    };

    // Debounced search effect
    useEffect(() => {
        const t = setTimeout(() => {
            if (searchValue === lastAppliedQRef.current) return;

            setSearchParams({ page: "1", limit: String(pagination.pageSize), q: searchValue || "" }, { replace: true });
            lastAppliedQRef.current = searchValue;
            fetchData(1, pagination.pageSize, searchValue);
        }, 400);

        return () => clearTimeout(t);
    }, [searchValue]);

    useEffect(() => {
        setTitle("Projects")
    })

    const columns = [
        {
            title: "Name",
            dataIndex: "name",
            key: "name",
            render: (project) => <a onClick={() => navigateToReleases(project)}>{project}</a>,
        },
        { title: "Description", dataIndex: "description", key: "description" },
        {
            title: "Created",
            dataIndex: "created_at",
            key: "created",
            render: (text) => <>{new Date(text).toLocaleString()}</>,
        },
    ];

    const navigateToReleases = (project) => {
        navigate(`/projects/${project}/releases`);
    };

    return (
        <>
            <Flex justify="space-between" style={{ marginBottom: 16 }}>
                <Input
                    placeholder="Search"
                    value={searchValue}
                    onChange={(e) => setSearchValue(e.target.value)}
                    style={{ marginRight: 16 }}
                />
                <Button type="primary" icon={<PlusOutlined />} onClick={() => setIsModalOpen(true)}>
                    New
                </Button>
            </Flex>

            <Table
                columns={columns}
                dataSource={data}
                rowKey="name"
                loading={loading}
                pagination={pagination}
                onChange={handleTableChange}
            />

            <Modal
                title="Create new project"
                open={isModalOpen}
                onCancel={() => setIsModalOpen(false)}
                footer={[
                    <Button form="createProject" type="primary" key="submit" htmlType="submit">
                        Create
                    </Button>,
                ]}
            >
                <div style={{ marginTop: 20 }}>
                    <CreateProject onSuccess={() => {
                        fetchData(pagination.current, pagination.pageSize, lastAppliedQRef.current);
                    }} />
                </div>
            </Modal>
        </>
    );
};

export default ProjectsList;