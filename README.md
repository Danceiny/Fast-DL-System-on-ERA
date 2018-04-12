# 毕设

***Graduation design of undergraduate***

Title: A Fast Deep-learning system based on Economic Resource Allocation for Cloud

基于ERA的快速深度学习系统。

改题目：

**面向深度学习的云端资源调度系统的设计与构建**

有关ERA，参见[Paper ERA](./Translations/ERA.md)

| Proceeding                               |
| ---------------------------------------- |
| [WWW '17 Companion](http://www.www2017.com.au/) Proceedings of the 26th International Conference on World Wide Web Companion |
| Pages 635-642                            |

图书馆文献查不到。。。只能找到两届WWW大会的，比较久远了。

School of Electronic Information and Communications, HUST, Spring 2018

## 运行
```
$ git clone https://github.com/Danceiny/Fast-DL-System-on-ERA.git {YOUR_DIR}

$ echo -e "if [ -z "$GOPATH" ]; then\n\texport GOPATH={YOUR_DIR}\nelse\n\texport GOPATH=$GOPATH:{YOUR_DIR}" >> ~/.zshrc && source ~/.zshrc # or ~/.bashrc...

$ cd {YOUR_DIR} && govendor add +external
```


### 开题答辩

1. 到底解决什么问题？

   主要解决深度学习平台（云端系统）的资源分配效率低下的问题。

2. 评判标准是什么？

   - 任务启动速度

   - - 各阶段启动速度
     - 总体启动速度 XX秒以内

   - 任务启动成功率

   - - 启动失败可能情况：系统bug；系统拒绝提供服务
     - 一次成功率、二次成功率等

   - 任务运行成功率

   - - 运行失败可能情况：系统bug

   - 云资源利用率

   - - 采用ERA方案的利用率与未采用ERA方案的利用率的比值，该比值越大越“成功”

   - 用户成本

   - - 以单位成本作为标准，用户实际消费金额/实际运行时间（ 同机型）
     - 采用ERA方案的单位成本与未采用ERA方案的单位成本的比值，该比值越小越“成功”

3. 你的工作是什么？

   根据“资源预定”的经济模型，开发相应的任务调度、动态定价、需求预测算法，作为核心调度层，并整合其他模块，构成面向深度学习的云端系统。

4. 与深度学习有什么关系？

   所设计和构建的云端系统，其应用领域是深度学习，其终端用户是深度学习研发者。


## 参考文章
- https://www.nextplatform.com/2017/03/02/experimental-cloud-reservation-agent-gets-right-work-done-faster/


## 日程备忘录
1. Go模块的数据建模，TCP服务器的处理流程梳理；
2. Python客户端user/job/jobreq/resource三个重要模型的建立；


## 调试
1. 启动celery

`sudo /Users/huangzhen/dev/miniconda3/envs/py27/bin/celery -A Platform.ERACenter.Cloud_Interface.cloud worker --autoscale=50,3 --loglevel=debug`