import { useEffect, useRef, useState} from "react";
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
    notification,
} from "antd";
import { LockOutlined, SearchOutlined, ReloadOutlined} from "@ant-design/icons";
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

const ConfigRow = ({ title, description, control, locked = false, query, isChanged }) => (
    <Row align="middle" justify="flex-start" style={{ padding: "12px 0" }}>
        <Col flex="auto" style={{ maxWidth: 500 }}>
            <Text strong>
                <Highlight text={title} query={query} />
            </Text>
            <br />
            <Text type="secondary">{description}</Text>
        </Col>
        <Col style={{ display: "flex", alignItems: "center", gap: 8 }}>
            <div style={isChanged ? { border: '2px solid #1890ff', borderRadius: '4px' } : {}}>
                {control}
            </div>
            {locked && <LockOutlined style={{ color: "#999" }} />}
        </Col>
    </Row>
);

const ConfigsList = () => {
    const {project, environment, release} = useParams();
    const groupRefs = useRef({});
    const [api, contextHolder] = notification.useNotification();
    const [configs, setConfigs] = useState([]);
    const [currentEnv, setCurrentEnv] = useState(environment)
    const [originalValues, setOriginalValues] = useState(new Map());
    const [modifiedValues, setModifiedValues] = useState(new Map());
    const [highlighted, setHighlighted] = useState(null);
    const [filter, setFilter] = useState("");

    useEffect(() => {
        fetchConfigs(currentEnv)
    }, []);

    const fetchConfigs = async (env) => {
        fetch(`/api/v1/projects/${project}/envs/${env}/releases/${release}/configs`)
            .then((response) => {
                if (!response.ok) throw new Error("err");
                return response.json();
            })
            .then((data) => {
                const configsData = data.data.configs;
                setConfigs(configsData);

                // Сохраняем оригинальные значения
                const newOriginalValues = new Map();
                configsData.forEach(cfg => {
                    newOriginalValues.set(cfg.key, cfg.value);
                });
                setOriginalValues(newOriginalValues);

                // Сбрасываем измененные значения
                setModifiedValues(new Map());
            });

        setCurrentEnv(env)
    }

    const handleChange = (key, value) => {
        const newValue = value?.toString() ?? "";
        const originalValue = originalValues.get(key);

        setConfigs((prev) => {
            return prev.map((cfg) => {
                if (cfg.key === key) {
                    return {
                        ...cfg,
                        value: newValue
                    };
                }
                return cfg;
            });
        });

        setModifiedValues((prev) => {
            const newMap = new Map(prev);

            if (newValue === originalValue) {
                // Если значение вернулось к оригинальному, удаляем из изменений
                newMap.delete(key);
            } else {
                // Сохраняем изменение
                newMap.set(key, {
                    oldValue: originalValue,
                    newValue: newValue
                });
            }

            return newMap;
        });
    };

    const handleSave = () => {
        // Создаем объект только с измененными значениями
        const changes = {};
        for (let [key, value] of modifiedValues) {
            changes[key] = value.newValue;
        }

        if (changes.length === 0) {
            return
        }

        fetch(`/api/v1/projects/${project}/envs/${currentEnv}/releases/${release}/configs`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(changes),
        })
            .then((response) => {
                if (!response.ok) {
                    api.error({
                        message: "Failed to update",
                        description: response.status
                    })
                    return
                }

                api.success({
                    message: "Updated success"
                })

                return fetchConfigs(currentEnv);
            })
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

    const handleReload = () => {
        fetchConfigs(currentEnv)
        api.success({
            message: "Updated"
        })
    }

    return (
        <>
            {contextHolder}
            <Button type="dashed" icon={<ReloadOutlined />} style={{marginRight: 20}} onClick={() => handleReload()}></Button>
            <EnvironmentsSelector
                project={project}
                environment={currentEnv}
                style={{marginBottom: 20}}
                onEnvChange={(env) => fetchConfigs(env)}
            />
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
                                        isChanged={modifiedValues.has(cfg.key)}
                                    />
                                    {idx < items.length - 1 && <Divider />}
                                </div>
                            ))}
                        </Card>
                    ))}
                    <Button
                        type="primary"
                        onClick={handleSave}
                        style={{ marginTop: 24 }}
                        disabled={modifiedValues.size === 0}
                    >
                        Save ({modifiedValues.size})
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