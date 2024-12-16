# Hessian V1 序列化数据解析

解析 Hessian V1 序列化数据

## Build

```
go build -o hs_serial_parse ./cmd/main.go
```

## Run

```
./hs_serial_parse --path examples/{hs_serial_path}.ser
```

## Tests

| Gadget      | Package | Parse |
| ----------- | ----------- | ----------- |
| SpringAbstractBeanFactoryPointcutAdvisor.ser      | marshalsec     | √ |
| SpringPartiallyComparableAdvisorHolder.ser        | marshalsec     | √ |
| Resin.ser                                         | marshalsec     | √ |
| Rome.ser                                          | marshalsec     | √ |
| XBean.ser                                         | marshalsec     | √ |

## Example

run:
```
./hs_serial_parse --path examples/SpringAbstractBeanFactoryPointcutAdvisor.ser
```

result:
```
{
  "ClassName": "",
  "Maps": [
    {
      "Key": {
        "ClassName": "org.springframework.aop.support.DefaultBeanFactoryPointcutAdvisor",
        "Maps": [
          {
            "Key": "adviceBeanName",
            "Value": "ldap://127.0.0.1"
          },
          {
            "Key": "order",
            "Value": null
          },
          {
            "Key": "pointcut",
            "Value": {
              "ClassName": "org.springframework.aop.TruePointcut",
              "Maps": []
            }
          },
          {
            "Key": "beanFactory",
            "Value": {
              "ClassName": "org.springframework.jndi.support.SimpleJndiBeanFactory",
              "Maps": [
                {
                  "Key": "resourceRef",
                  "Value": true
                },
                {
                  "Key": "shareableResources",
                  "Value": {
                    "ListType": "java.util.HashSet",
                    "Args": [
                      "ldap://127.0.0.1"
                    ]
                  }
                },
                {
                  "Key": "singletonObjects",
                  "Value": {
                    "ClassName": "",
                    "Maps": []
                  }
                },
                {
                  "Key": "resourceTypes",
                  "Value": {
                    "ClassName": "",
                    "Maps": []
                  }
                },
                {
                  "Key": "logger",
                  "Value": {
                    "ClassName": "org.apache.commons.logging.impl.NoOpLog",
                    "Maps": []
                  }
                },
                {
                  "Key": "jndiTemplate",
                  "Value": {
                    "ClassName": "org.springframework.jndi.JndiTemplate",
                    "Maps": [
                      {
                        "Key": "logger",
                        "Value": {
                          "ClassName": "org.apache.commons.logging.impl.NoOpLog",
                          "Maps": []
                        }
                      },
                      {
                        "Key": "environment",
                        "Value": null
                      }
                    ]
                  }
                }
              ]
            }
          }
        ]
      },
      "Value": null
    },
    {
      "Key": null,
      "Value": null
    },
    {
      "Key": null,
      "Value": null
    },
    {
      "Key": {
        "ClassName": "org.springframework.aop.support.DefaultBeanFactoryPointcutAdvisor",
        "Maps": [
          {
            "Key": "adviceBeanName",
            "Value": null
          },
          {
            "Key": "order",
            "Value": null
          },
          {
            "Key": "pointcut",
            "Value": null
          },
          {
            "Key": null,
            "Value": null
          },
          {
            "Key": null,
            "Value": null
          },
          {
            "Key": "beanFactory",
            "Value": null
          }
        ]
      },
      "Value": null
    },
    {
      "Key": null,
      "Value": null
    },
    {
      "Key": null,
      "Value": null
    }
  ]
}
```