import { useEffect, useRef, useState } from "react";
import { useParams } from 'react-router-dom';
import {
    Card,
    Button,
    Input,
    InputNumber,
    Switch,
    Typography,
    Row,
    Col,
    Divider,
} from "antd";
import { LockOutlined, SearchOutlined } from "@ant-design/icons";
import EnvironmentsSelector from "../Environments/Selector";

const { Text } = Typography;

const Highlight = ({ text, query }) => {
    if (!query) return <>{text}</>;
    const regex = new RegExp(`(${query})`, "gi");
    const parts = text.split(regex);
    return (
        <>
            {parts.map((part, i) =>
                part.toLowerCase() === query.toLowerCase() ? (
                    <span key={i} style={{ backgroundColor: "yellow" }}>
                        {part}
                    </span>
                ) : (
                    part
                )
            )}
        </>
    );
};

const ConfigRow = ({ title, description, control, locked = false, query }) => (
    <Row align="middle" justify="flex-start" style={{ padding: "12px 0" }}>
        <Col flex="auto" style={{ maxWidth: 500 }}>
            <Text strong>
                <Highlight text={title} query={query} />
            </Text>
            <br />
            <Text type="secondary">{description}</Text>
        </Col>
        <Col style={{ display: "flex", alignItems: "center", gap: 8 }}>
            {control}
            {locked && <LockOutlined style={{ color: "#999" }} />}
        </Col>
    </Row>
);

const ConfigsList = () => {
    const {project, environment, release} = useParams();
    const [configs, setConfigs] = useState([]);
    const [highlighted, setHighlighted] = useState(null);
    const [filter, setFilter] = useState("");
    const groupRefs = useRef({});

    useEffect(() => {
        fetchConfigs(environment)
    }, []);

    const fetchConfigs = async (env) => {
        fetch(`/api/v1/projects/${project}/envs/${env}/releases/${release}/configs`)
            .then((response) => {
                if (!response.ok) throw new Error("err");
                return response.json();
            })
            .then((data) => {
                setConfigs(data.data.configs);
            });
    }

    const handleChange = (key, value) => {
        setConfigs((prev) =>
            prev.map((cfg) =>
                cfg.key === key ? { ...cfg, value: value?.toString() ?? "" } : cfg
            )
        );
    };

    const handleSave = () => {
        console.log("SAVE:", configs);
    };

    // группировка
    const grouped = configs.reduce((acc, cfg) => {
        const g = cfg.group.toUpperCase() || "UNSPECIFIED";
        if (!acc[g]) acc[g] = [];
        acc[g].push(cfg);
        return acc;
    }, {});

    // фильтр
    const filteredGroups = Object.entries(grouped).reduce((acc, [group, items]) => {
        if (!filter) {
            acc[group] = items;
        } else {
            const matched = items.filter((cfg) =>
                cfg.key.toLowerCase().includes(filter.toLowerCase())
            );
            if (matched.length > 0) acc[group] = matched;
        }
        return acc;
    }, {});

    const renderControl = (cfg) => {
        const commonProps = {
            disabled: !cfg.writable,
            style: { minWidth: 120 },
        };

        switch (cfg.value_type) {
            case "string":
                return (
                    <Input
                        {...commonProps}
                        value={cfg.value}
                        onChange={(e) => handleChange(cfg.key, e.target.value)}
                    />
                );
            case "bool":
                return (
                    <Switch
                        {...commonProps}
                        checked={cfg.value === "true"}
                        onChange={(val) => handleChange(cfg.key, val)}
                    />
                );
            case "int":
            case "int64":
            case "uint":
            case "uint64":
            case "float":
            case "float64":
                return (
                    <InputNumber
                        {...commonProps}
                        value={Number(cfg.value)}
                        onChange={(val) => handleChange(cfg.key, val)}
                    />
                );
            default:
                return <Text type="secondary">Unsupported type: {cfg.value_type}</Text>;
        }
    };

    const handleScrollToGroup = (group) => {
        const el = groupRefs.current[group];
        if (el) {
            el.scrollIntoView({ behavior: "smooth", block: "start" });
            setHighlighted(group);
            setTimeout(() => setHighlighted(null), 1500);
        }
    };

    return (
        <>
            <EnvironmentsSelector project={project} environment={environment} style={{marginBottom: 20}} onEnvChange={(env) => fetchConfigs(env)} />
            <Row gutter={24} style={{ alignItems: "flex-start" }}>
                <Col span={20}>
                    <Input
                        placeholder="filter config"
                        prefix={<SearchOutlined />}
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        style={{ marginBottom: 16 }}
                    />
                    {Object.entries(filteredGroups).map(([group, items]) => (
                        <Card
                            key={group}
                            ref={(el) => (groupRefs.current[group] = el)}
                            id={group}
                            title={
                                <span
                                    style={{
                                        backgroundColor:
                                            highlighted === group ? "yellow" : "transparent",
                                        transition: "background-color 0.5s ease",
                                    }}
                                >
                                    {group}
                                </span>
                            }
                            style={{ marginBottom: 24 }}
                        >
                            {items.map((cfg, idx) => (
                                <div key={cfg.key}>
                                    <ConfigRow
                                        title={cfg.key}
                                        description={cfg.usage}
                                        control={renderControl(cfg)}
                                        locked={!cfg.writable}
                                        query={filter}
                                    />
                                    {idx < items.length - 1 && <Divider />}
                                </div>
                            ))}
                        </Card>
                    ))}
                    <Button type="primary" onClick={handleSave} style={{ marginTop: 24 }}>
                        Сохранить
                    </Button>
                </Col>
                <Col span={4}>
                    <Card title="Groups">
                        <div style={{ position: "sticky", top: 20 }}>
                            <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
                                {Object.keys(grouped).map((group) => (
                                    <a
                                        key={group}
                                        onClick={() => handleScrollToGroup(group)}
                                    >
                                        {group.toLowerCase()}
                                    </a>
                                ))}
                            </div>
                        </div>
                    </Card>
                </Col>
            </Row>
        </>
    );
};

export default ConfigsList;
