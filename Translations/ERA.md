# ERA - A Framework for Economic Resource Allocation for the Cloud (ERA - 为云计算而生的经济型资源分配框架）
![](http://opkk27k9n.bkt.clouddn.com/18-3-3/81518610.jpg)
>总字符数：	72119
>总字符数(不含空白)：	63001
>空白字符数：	9118
>中文字符数：	12216
>英文字符数：	46264
>标点符号数：	3500
>其它字符数：	1021
>中英文单词数：	22121
>非中文单词数：	9905
>内容行数：	477

## Abstract 摘要
Cloud computing has reached significant maturity from a systems perspective, but currently deployed solutions rely on rather basic economics mechanisms that yield suboptimal allocation of the costly hardware resources. In this paper we present Economic Resource Allocation (ERA), a complete framework for scheduling and pricing cloud resources, aimed at increasing the efficiency of cloud resources usage by allocating resources according to economic principles. The ERA architecture carefully abstracts the underlying cloud infrastructure, enabling the development of scheduling and pricing algorithms independently of the concrete lower-level cloud infrastructure and independently of its concerns. Specifically, ERA is designed as a flexible layer that can sit on top of any cloud system and interfaces with both the cloud resource manager and with the users who reserve resources to run their jobs. The jobs are scheduled based on prices that are dynamically calculated according to the predicted demand. Additionally, ERA provides a key internal API to pluggable algorithmic modules that include scheduling, pricing and demand prediction. We provide a proof-of-concept software and demonstrate the effectiveness of the architecture by testing ERA over both public and private cloud systems –Azure Batch of Microsoft and Hadoop/YARN. A broader intent of our work is to foster collaborations between economics and system communities. To that end, we have developed a simulation platform via which economics and system experts can test their algorithmic implementations.

从系统的角度来看，云计算已经达到了显著的成熟度，但目前部署的解决方案依赖于基本的经济机制，这些机制产生了昂贵的硬件资源的次优分配。在本文中，我们提出了经济资源分配（ERA），一个完整的框架，用于调度和定价云资源，旨在通过根据经济原则分配资源，提高云资源使用的效率。该架构小心地抽象了底层的云基础设施，使开发调度和定价算法不依赖于具体的底层云基础设施，而不依赖于其关注点。具体来说，ERA被设计成一个灵活的"层"，可以配置于任何云系统的顶部，同时与云资源管理器和预留资源的用户一起运行他们的工作。作业是根据根据预期需求动态计算的价格排定的。此外，ERA为可插入算法模块提供了一个关键的内部API，包括调度、定价和需求预测。我们提供了一个概念软件的证明，并通过测试ERA部署在公有云和私有云系统——Azure Batch of Microsoft和Hadoop/YARN，证明了该架构的有效性。我们工作的一个更广泛的意图是促进经济学和系统社区之间的协作。为此，我们开发了一个仿真平台，通过这个平台，经济学和系统专家可以测试他们的算法实现。

## 关键词
Cloud Computing; Economics; Dynamic Pricing; Reservations

云计算；经济学；动态计价；预留

## 1. 介绍
Cloud computing, in its private or public incarnations, is commonplace in industry as a paradigm to achieve high return on investments (ROI) by sharing massive computing infrastructures that are costly to build and operate [5, 14]. Effective sharing pivots around two key ingredients: 1) a system infrastructure that can securely and efficiently multiplex a shared set of hardware resources among several tenants, and 2) economic mechanisms to arbitrate between conflicting resource demands from multiple tenants. State-of-the-art cloud offerings provide solutions to both, but with varying degrees of sophistication. The system challenge has been subject to extensive research focusing on space and time multiplexing. Space multiplexing consists of sharing servers among tenants, while securing them via virtual machine [34, 32, 7] and container technologies [23, 22]. Time multiplexing comprises a collection of techniques and systems that schedule tasks over time. The focus ranges from strict support of Service Level Objectives (SLOs) [10, 17, 30, 11] to maximization of cluster utilization [12, 35, 18, 25, 24, 13, 8]. Many of these advances are already deployed solutions in the public cloud [26, 4] and the private cloud [1, 31, 15, 33, 8]. This indicates a good degree of maturity in how the system challenge is tackled in cloud settings. On the other hand, while the economics challenge has received some attention in recent research (see, e.g., [27, 16, 6, 28, 19, 20] and references therein), the underlying principles have not yet been translated into equally capable solutions deployed in practice.

云计算，不管是私有云还是公有云这些衍生版本,通过共享巨大的计算基础设施——而建造和操作这些基础设施都是昂贵的，都是工业界的实现高投资回报比（ROI）的范例[ 5, 14 ]。高效的共享中枢围绕着两个关键要素：1）能够安全高效地在用户间复用一系列共享硬件资源的一套系统基础设施；2）一套在多个用户的资源需求冲突时进行仲裁的经济机制。

最先进的云计算提供上述两方面的解决方案，但都有不同程度的复杂性。该系统的挑战有这样的倾向：聚焦在空间和时间复用的扩展性研究。空间复用包括用户间共享服务器，其安全性由虚拟机以及容器技术保证。时间复用包含一系列的调度任务的技巧和系统。该系统的挑战一直受到广泛的研究，重点是空间和时间复用。空间复用包括在租户之间共享服务器，同时通过虚拟机[ 34, 32, 7 ]和容器技术【23, 22】来保护它们。时间多路复用包括一系列技术和系统，它们可以根据时间安排任务。重点从服务水平目标（SLO）严格支持[ 10, 17, 30，11 ]最大化集群利用[ 12, 35, 18，25, 24, 13，8 ]。其中许多进展已经在公有云[ 26, 4 ]和私有云中部署了解决方案。[ 1, 31, 15，33, 8 ]。这表明在云环境中如何解决系统挑战具有一个很好的成熟度。另一方面，虽然经济挑战在最近的研究中得到了一些关注（例如，[ 27, 16, 6，28, 19, 20 ]和其中的参考文献），但基本原则尚未转化为在实践中部署的同样有能力的解决方案。

### 1.1 经济挑战和ERA的方法
In current cloud environments, resource allocation is governed by very basic economics mechanisms. The first type of mechanism (common in private clouds [31, 8, 15, 33]) uses fixed pre-paid guaranteed quotas. The second type (common in public clouds [26, 4]) uses on-demand unit prices: the users are charged real money（Or, within a company, fiat money） per unit of resource used. In most cases these are fixed prices, with the notable exception of Amazon’s spot instances that use dynamically changing prices.（Although, based on independent analysis, even these may not truly leverage market mechanisms to determine prices [3].） Spot-instance offerings, however, do not provide guaranteed service, as the instances might be evicted if the user bid is too low. Hence, utilizing spot instances requires special attention from the user when determining his bid, and might not be suitable for high-priority production jobs [21, 28, 2]. The fundamental problem is finding a pricing and a scheduling scheme that will result in highly desired outcome, that of high efficiency.

在当前的云环境中，资源分配受非常基本的经济学机制支配。第一种机制（常见于私有云（31, 8, 15，33））使用固定的预先支付的保证配额。第二种（公共云中常见的26, 4种）使用按需单价：用户每单位所使用的资源是实实在在的钱（或在公司内，法定货币）。在大多数情况下，这些都是固定的价格，使用动态价格的亚马逊的spot instances（现货实例）例外。（虽然，基于独立分析，甚至不可能真正利用市场机制确定价格[ 3 ]。）spot instances的产品，但不提供保证服务的情况下，可能会拆迁户如果用户出价太低。因此，在确定出价时使用spot instances需要特别注意，并且可能不适合于高优先级生产作业[ 21, 28, 2 ]。最根本的问题是找到一个定价和一个调度方案，这将导致高期望的结果，高效率。

***Efficiency***: From an economics point of view, the most fundamental goal for a cloud system is to maximize the economic efficiency, that is, to maximize the total value that all users get from the system. For example, whenever two users have conflicting demands, the one with the lowest cost for switching to an alternative (running at a different time/place or not running at all) should be the one switching. The resources would thus be allocated to the user with “highest marginal value.” The optimal-allocation benchmark for a given cloud is that of an omniscient scheduler who has access to the complete information of all cloud users -——including their internal costs and alternative options —— and decides what resources to allocate to whom in a way that maximizes the efficiency goals of the owner of the cloud. Let us stress: to get a meaningful measure of efficiency we must count the value obtained rather than the resources used, and we should aim to maximize this value-based notion of efficiency.（Another important goal, of course, is revenue, but we note that the potential revenue is limited by the value created, so high efficiency is needed for high revenue. Moreover, since there is competition between cloud providers, these providers generally aim to increase short-term efficiency as this is likely to have positive longterm revenue effects. The issue of increasing revenue is usually attacked under the “Platform as a Service” (PaaS) strategy of providing higher-level services. This is essentially orthogonal to allocation efficiency and is beyond the scope of the present paper.）

***效率***：从经济学的角度来看，云系统最基本的目标是最大化经济效率，即最大化所有用户从系统中获得的总价值。例如，当两个用户有冲突的需求时，切换到另一个（在不同时间/地点运行或根本不运行）的成本最低的，应该在冲突中切换备选方案。资源将被分配到“最高边际价值用户”。对于一个给定的云优化配置的基准是一个无所不知的调度器，它拥有至高的权限，包括所有的云用户的完整信息——包括内部成本和替代方案——和决定资源分配给谁，以这种方式使云的所有者的效率目标最大化。让我们强调：要获得有意义的测量效率，我们必须计算获得的价值，而不是使用的资源，而且我们应该最大化这个价值基础概念的效率。（另一个重要的目标，当然，是收入，但我们注意到，潜在的收益受限于创造的价值，要想获得高收入，效率高是必要的。此外，由于云提供商之间存在竞争，这些提供商通常着眼于提高短期效率，因为这很可能产生正面的长期收入效应。增加收入的问题通常 在提供高层次的服务的“平台即服务”（PaaS）策略下被攻击。这与分配效率基本上是正交的，超出了本论文的范围。）

***Limitations of current solutions***: With this value-based notion
of efficiency in mind, let us evaluate commonly deployed pricing mechanisms. Private cloud frameworks [31, 8, 15, 33] typically resort to pre-paid guaranteed quotas. The main problem with this option is that, formally, it provides no real sharing of common resources: to guarantee that every user always has his guaranteed pre-paid resources available, the cloud system must actually hold sufficient resources to satisfy the sum of all promised capacities, even though only a fraction will likely be used at any given time. Mechanisms such as work-preserving fair-queueing [12] are typically designed to increase utilization [31, 15], but they do not fundamentally change the equation for value-based efficiency, as the resources offered above the user quota are typically distributed at no cost and with no guarantees. Furthermore, lump-sum prepayment implies that the users’ marginal cost for using their guaranteed resources is essentially zero, and so they will tend to use their capacity for “non-useful” jobs whenever they do not fill their capacity with “useful” jobs. This often results in cloud systems that seem to be operating at full capacity from every engineering point of view, but are really working at very “low capacity” from an economics point of view, as most of the time, most of the jobs have very low value.
On the other hand, public cloud offerings [26, 4] typically employ on-demand unit-pricing schemes. The issue with this solution is that the availability of resources cannot be guaranteed in advance. Typically the demand is quite spiky, with short periods of peak demand interspersed within much longer periods of low demand. Such spiky demand is also typical of many other types of shared infrastructure such as computer network bandwidth, electricity, or ski resorts. In all these cases the provider of the shared resource faces a dilemma between extremely expensive over-provisioning of capacity and giving up on guaranteed service at peak times. In the case of cloud systems, for jobs that are important enough, users cannot take the risk of their jobs not being able to run when needed, and thus “on-demand” pricing is only appealing to low-value or highly time-flexible jobs, while most important “production” jobs with little time flexibility resort to buying long-term guaranteed access to resources. While flexible unit prices (such as Amazon’s spot instances) may have an advantage over fixed ones as they can better smooth demand over time, as highlighted above, they only get an opportunity to do so for the typically low-value “non-production” jobs.

***当前解决方案的局限性***：秉持着基于价值的效率理念，让我们计算通常部署的计价机制。私有云框架[ 31, 8, 15，33 ]，通常采用预先支付的保障性配额。这种方案的主要问题是，它没有提供真正的公用资源共享：为了保证每一位用户总是能获得可用的预付费的保障性资源，云系统实际上必须保持充足的资源来满足所有允诺载力之和，尽管在任何时刻只有一部分可能被使用。像这种资源预留-公平队列的机制，通常为增加利用率而设计，但是他们没有在根本上改变基于价值的效率等式，因为超过用户配额的提供资源通常没有分布式开销，也没有保障。一次性预付说明用户使用保障性资源的边际花费本质上是0，因此他们倾向于使用他们的容载来做一些无用的作业，不管他们有没有做有用的作业来充分利用容载。这经常导致云系统似乎从所有的工程学角度看都是在满载运行，尽管从经济学角度看是低负载的运行态，因为大部分时间，大部分的作业都是低价值的。

另一方面，公有云（26, 4）通常使用按需定价方式。这种解决方案的问题是，资源的可用性无法提前保障。通常需求是相当尖刻的，短时高频需求散布在长时间低频需求之中。这些需求的情况和其他的一些共享基础设施很像，例如计算机网络带宽，电力系统，滑雪胜地等。所有这些案例中，共享资源的提供方在极昂贵的容载冗余供应和放弃高峰期的保障性服务这二者间面临一个两难情境。在这些云系统中，对于足够重要的作业，用户不能承担作业不能被运行的风险。因此按需定价只适合低价值或者时间灵活不敏感的作业，尽管大多数重要的时间敏感的生产作业倾向于购买长期的保障性资源访问权。尽管灵活的单价（例如亚马逊的现货实例）可能比固定单价有优势，因为前者可以在时间上平滑需求。如前文所强调，它们仅仅有机会服务于低价值的非生产作业。

***The ERA approach***: The pricing model that we present in ERA enables sharing of resources and smoothing of demand even for high-value production jobs. This is done using the well-known notion of reservations, commonly used for many types of resources such as hotel rooms or super-computer time as well as in a few cloud systems [10, 17, 30], but in a flexible way in terms of both pricing and scheduling. We focus on the economic challenges of scheduling and pricing batch style computations with completion time SLOs (deadlines) on a shared cloud infrastructure. The basic model presented to the user is that of resource reservation. At “reservation time,” the user’s program specifies its reservation request. The basic form of such a request is:（The general form of requests is given by ERA’s “bidding description language” (see Section 2.3.1) that allows specifying multiple resources, variable “shapes” of use across time, and combinations thereof.）
>Basic Reservation: “I need 100 containers (with 6GB and 2cores each) for 5 hours, some time between 6am and 6pm today, and am willing to pay up to $100 for it.”

This class of workloads is very prominent – much of “big data” falls under this category [11, 30, 17] – and it provides us with a unique opportunity for economic investigation. While state-of-theart solutions provide effective system-level mechanisms for sharing resources, they rely on users’ goodwill to truthfully declare their resource needs and deadlines to the system. By dynamically manipulating the price of resources, ERA provides users with incentives to expose to the system as much flexibility as possible. The ERA mechanism ensures that the final price paid by the user is the lowest possible price the system can provide for granting the request. The more flexibility a user exposes, the better the user’s chances of getting a good price. If this minimal price exceeds the stated maximal willingness to pay, then the request is denied.（Alternatively, the ERA mechanism may present this minimal price as a “quote” to the user, who may decide whether to accept or reject it） Once a reservation request is accepted, the payment is fixed at reservation time, and the user is assured that the reserved resources will be available to him within the requested window of time. The guarantee is to satisfy the request rather than provide a promise of specific resources at specific times. For more details regarding the model presented to the user see Section 2.3.1.

***ERA的方法***：我们在ERA中提供的定价模式可以实现资源共享和平滑需求，即使对于高价值的生产工作也是如此。 这是通过使用众所周知的预订来完成的，通常用于许多类型的资源，例如酒店房间或超级计算机时间以及一些云系统[10,17,30]，但以灵活的方式定价和计划的条款。 我们专注于调度和定价批处理类型计算的经济挑战，并在共享云基础架构上使用完成时间SLO（截止日期）。呈现给用户的基本模型是资源预留。 在“预订时间”，用户的程序指定其预订请求。 这种请求的基本形式是：（一般形式的请求由ERA的“招标描述语言”（见第2.3.1节）给出，允许指定多个资源，随时间变化的“形状”，及其组合。）
>基本预订：我需要100个容器（每个6GB的内存和2CPU核心），5个小时，今天6：00到18：00之间的某个时间段，我可以付给你100美元。

这类工作量非常突出 - 大部分“大数据”属于这个类别[11,30,17] - 它为我们提供了一个独特的经济调查机会。尽管现有的解决方案为共享资源提供了有效的系统级机制，但它们依赖用户的善意——将真实的资源需求和截止日期提交给系统。通过动态地操纵资源价格，ERA鼓励用户向系统展示尽可能多的灵活性。ERA机制确保用户支付的最终价格是系统为授予请求所能提供的最低价格。用户暴露的灵活性越高，用户获得理想价格的机会就越大。如果这个最低价格超过了用户的最大支付意愿，那么请求被拒绝（或者，ERA机制可以将这个最低价格作为“报价”呈现给用户，用户可以决定是接受还是拒绝）。一旦预约请求被接受，支付在预定时间被固定，并且用户被确信在所请求的时间窗口内预留的资源将对他可用，保证是满足要求，而不是在特定时间提供具体资源的承诺。有关向用户提供的模型的更多详细信息，请参见第2.3.1节。

### 1.2 ERA总览
A key part of the challenge of devising good allocation schemes for cloud resources is their multi-faceted nature: while our goals are in terms of economic high-level business considerations, implementation must be directly carried out at the computer systemsengineering level. These two extreme points of view must be connected using clever algorithms and implemented using appropriate software engineering. Indeed, in the literature regarding cloud systems, one can see many papers that deal with each one of these aspects – “systems” papers as well as theoretical papers on scheduling and pricing jobs in cloud systems – as cited above. Unfortunately, these different threads in the literature are often disconnected from each other and they are not easy to combine to get an overall solution. We believe that one key challenge at this point is to provide a common abstraction that encompasses all these considerations. We call this the architectural challenge.

制定云资源良好分配方案的挑战中的关键部分是它们的多面性：虽然我们的目标是依据经济的高级别业务所进行的考量，但实施必须直接在计算机系统工程级别执行。这两个极端的观点必须使用巧妙的算法进行连接，并使用适当的软件工程来实现。 事实上，在关于云系统的文献中，如上所述，可以看到许多论文处理这些方面中的每一个——“系统”论文以及关于云系统中调度和定价作业的理论论文。不幸的是，文献中的这些不同线索常常彼此断开，并且它们不易组合以获得整体解决方案。我们认为在这一点上的一个关键挑战是提供一个包含所有这些考虑因素的共同抽象。我们称之为架构挑战。

![](http://opkk27k9n.bkt.clouddn.com/18-3-3/17436991.jpg)
Figure 1: ERA Architecture. The ERA system is designed as an intermediate layer between the users and the underlying cloud infrastructure. The same actual core code is also interfaced with the simulator components.
图1：ERA架构 ERA系统被设计为用户和底层云基础设施之间的中间层。 相同的实际核心代码也与模拟器组件接口。

![](http://opkk27k9n.bkt.clouddn.com/18-3-3/98577912.jpg)
Figure 2: ERA Simulator Screenshot

图2：ERA模拟器截屏

Our answer to this challenge is the ERA system (Economic Resource Allocation). The ERA system is designed as an intermediate layer between the cloud users and the underlying cloud infrastructure. It provides a single model that encompasses all the very different key issues underlying cloud systems: economic, algorithmic, systems-level, and human-interface ones. It is designed to integrate economics insights practically in real-world system infrastructures, by guiding the resource allocation decisions of a cloud system in an economically principled way.（Not all resources in the cloud have to be managed via ERA. It is also possible that the cloud will let ERA manage only a subset of the resources (allowing the system to be incrementally tested), or will have several instances of ERA to manage different subsets of the resources.） This is achieved by means of a key architectural abstraction: the ERA Cloud Interface, which hides many of the details of the resource management infrastructure, allowing ERA to tackle the economic challenge almost independently of the underlying cloud infrastructure. ERA satisfies three key design goals: 1) it provides a crisp theoretical abstraction that enables more formal studies; 2) it is a practical end-to-end software system; and 3) it is designed for extensibility, where all the algorithms are by design easy to evolve or experiment with.

我们应对这一挑战的答案就是ERA系统。 ERA系统被设计为云用户和底层云基础设施之间的中间层。它提供了一个单一模型，涵盖了云系统所有非常不同的关键问题：经济，算法，系统级和人机界面。它旨在通过经济原则的方式指导云系统的资源分配决策，从而将经济学见解实际上集成到现实世界的系统基础设施中。（并非云中的所有资源都必须通过ERA进行管理。云也可能让ERA只管理资源的一个子集（允许系统进行增量测试），或者将有几个ERA实例来管理不同的资源子集。）这可以通过关键架构抽象达成：ERA云接口隐藏了资源管理基础架构的许多细节，使ERA几乎独立于底层云基础架构来应对经济挑战。 ERA满足三个关键设计目标：1）它提供了一个清晰的理论抽象，能够进行更正式的研究; 2）它是一个实用的端到端软件系统; 3）它是为扩展性而设计的，所有的算法都是通过设计易于进化或实验的。

***ERA’s key APIs***: ERA has two main outward-facing APIs as well as a key internal API. Figure 1 gives a high-level schematic of the architecture. The first external API faces the users and provides them with the economic reservation model of cloud services described above. The second external API faces the low-level cloud resource manager. It provides a separation of concerns that frees the underlying cloud system from any time-dependent scheduling or from any pricing concerns, and frees the ERA system from the burden of assigning specific processors to specific tasks in a reasonable resource-locality way, or from the low-level mechanics of firing up processes or swapping them out. See more details in Section 2.3.



***ERA的关键API***：ERA有两个主要的面向外部的API以及一个关键的内部API。图1给出了该架构的高级示意图。第一个外部API面向用户，并为他们提供上述云服务的经济预留模型。第二个外部API面向低级云资源管理器。它提供了一个关注点的分离，可以将基础云系统从任何时间相关的调度或任何定价问题中解放出来，并且可以使ERA系统免于以合理的资源区域方式将特定处理器分配给特定任务的负担，低级别的机制来启动进程或将其交换出去。更多细节见第2.3节。



Finally, the internal API is to pluggable algorithmic scheduling, pricing, and prediction modules. Our basic scheduling and pricing algorithm dynamically computes future resource prices based on supply and demand, where the demand includes both resources that are already committed to and predicted future requests, and schedules and prices the current request at the “cheapest” possibility. Our basic prediction model uses traces of previous runs to estimate future demand. The flexible algorithmic API then allows for future algorithmic, learning, and economic optimizations. The internal interfaces as well as our basic algorithmic implementations are described in Section 3.





最后，内部API是可插入算法调度、定价和预测模块。我们的基本调度和定价算法根据供需情况动态计算未来的资源价格，其中需求既包括已经承诺的预测未来需求的资源，也包括预测未来需求的资源，并以“最便宜”的可能性对当前需求进行计划和定价。我们的基本预测模型使用之前运行的痕迹来估计未来的需求。灵活的算法API允许未来的算法、学习和经济优化。第3节介绍了内部接口以及我们的基本算法实现。



Our goal in defining this abstraction is more ambitious than mere good software engineering in our system. As part of the goal of fostering a convergence between system and economic considerations, we have also built a flexible cloud simulation framework. The simulator provides an evaluation of key metrics, both “system ones” such as loads or latency, as well as “economic ones” such as “welfare” or revenue, as well as provides a visualization of the results (see screenshot in Figure 2). The simulator was designed to provide a convenient tool both for the cloud system’s manager who is interested in evaluating ERA’s performance as a step toward integration and for researchers who develop new algorithms for ERA and are interested in experimenting with their implementation without the need to run a large cluster. As is illustrated in Figure 1, the same core code that receives actual user requests and runs over the underlying cloud resource manager may be connected instead to the simulator so as to test it under variable loads and alternative cloud models. Comparing the results from our simulator and physical cluster runs, we find the simulator to be faithful (Section 4).

我们定义这种抽象的目标比我们的系统中纯粹的软件工程更雄心勃勃。作为促进系统和经济考虑之间融合的目标的一部分，我们还构建了灵活的云模拟框架。模拟器提供关键指标的评估，包括负载或延迟等“系统指标”以及“福利”或收入等“经济指标”，并提供结果可视化（请参见图2中的屏幕截图）。该模拟器旨在为云系统的管理人员提供便利的工具，该管理人员有兴趣评估ERA的性能，以此作为实现整合的一个步骤，以及为ERA开发新算法的研究人员，并且有兴趣在无需运行大集群。如图1所示，接收实际用户请求并通过底层云资源管理器运行的相同核心代码可以连接到模拟器，以便在可变负载和备选云模型下对其进行测试。比较我们模拟器和物理集群运行的结果，我们发现模拟器是可信的（第4节）。

The ERA system is implemented in Java, and an alternative implementation (of a subset of ERA) in C# was also done. We have performed extensive runs of ERA within the simulator as well as proof-of-concept runs with two prominent resource managers in the public and private clouds: the full system was interfaced with Hadoop/YARN [31] and the C# version of the code was interfaced and tested with Microsoft’s Azure Batch（Azure Batch is a cloud-scale job-scheduling and compute management service  https://azure.microsoft.com/en-us/services/batch/） simulator [29]. These runs show that the ERA algorithms succeed in increasing the efficiency of cloud usage, and that ERA can be successfully integrated with real cloud systems. Additionally, we show that the ERA simulator gives a good approximation to the actual run on a cloud system and thus can be a useful tool for developing and testing new algorithms. In Section 4 we present the results of a few of these runs.



ERA系统是用Java实现的，而C＃中的一个替代实现（ERA的一个子集）也已完成。 我们在模拟器中执行了广泛的ERA运行，并且在公有云和私有云中使用两个着名的资源管理器进行了概念证明运行：完整的系统与Hadoop / YARN [31]以及代码的C＃版本 与微软的Azure Batch（Azure批处理是一种云规模的作业调度和计算管理服务。https://azure.microsoft.com/en-us/services/batch/）模拟器进行了接口和测试[29]。 这些运行表明，ERA算法能够成功提高云使用效率，并且ERA可以与真正的云系统成功整合。 另外，我们证明ERA模拟器可以很好地近似于云系统上的实际运行，因此可以成为开发和测试新算法的有用工具。 在第4节中，我们将介绍这些运行的一些结果。

***Contributions***: In summary, we present ERA, a reservation system for pricing and scheduling jobs with completion-time SLOs. ERA makes the following contributions:

1. We propose an abstraction and a system architecture that allows us to tackle the economic challenge orthogonally to the underlying cloud infrastructure.
2. We devise algorithms for scheduling and pricing batch jobs with SLOs, and for predicting resource demands.
3. We design a faithful cloud simulator via which economics and system experts can study and test their algorithmic implementations.
4. We integrate ERA with two cloud infrastructures and demonstrate its effectiveness experimentally.

***贡献***：总而言之，我们提出了ERA，这是一个用于定价和安排具有完成时间SLO的工作的预订系统。 ERA做出以下贡献：
1. 我们提出抽象和系统架构，使我们能够正确处理与底层云基础架构相关的经济挑战。
2. 我们设计了用SLO调度和定价批量作业的算法，并预测了资源需求。
3. 我们设计了一个忠实的云模拟器，通过它可以让经济学和系统专家学习和测试他们的算法实现。
4. 我们将ERA与两个云基础架构集成在一起，并通过实验展示其有效性。


## 2. THE ERA MODEL AND ARCHITECTURE （ERA模型和架构）
### 2.1 The Bidding Reservation Model with Dynamic Prices （使用动态价格的配额预定模型）

ERA is designed to handle a set of computational resources of a cloud system, such as cores and memory, with the goal of allocating these resources to users efficiently. The basic idea is that a user that needs to run a job at some future point in time can make a reservation for the required resources and, once the reservation is accepted, these resources are then guaranteed (insofar as physically possible) to be available at the reserved time. The guarantee of availability of reserved resources allows high-value jobs to use cloud just like in the pre-paid guaranteed quotas model, but without the need to buy the whole capacity (for all times), which thus also allows for time sharing of resources, which increases efficiency.

The price for these resources is quoted at reservation time and is dynamically computed according to (algorithmically estimated) demand and the changing supply. More user flexibility in timing is automatically rewarded by lower prices. The basic idea is that these dynamic prices will regulate demand, achieving a better utilization of cloud resources. This mechanism also ensures that at peak times – where demand can simply not be met – the most “valuable” jobs are the ones that will be run rather than arbitrary ones. ERA uses a simple bidding model in which the user attaches to each of his job requests a monetary value specifying the maximal amount he is willing to pay for running the job. The amount of value lost for jobs that cannot be accommodated at these peak times serves as a quantification of the value that will be gained by buying additional cloud resources, and is an important input to the cloud provider’s purchasing decision process.

ERA旨在处理云系统的一组计算资源，如核心和内存，其目标是有效地将这些资源分配给用户。其基本思想是，需要在将来某个时间点运行作业的用户可以对所需资源进行预定，并且一旦接受预约，就可以保证这些资源（在物理上可能的情况下）在保留时间。预留资源可用性的保证允许高价值作业像预付费保证配额模型那样使用云，但无需购买整个容量（始终），因此也允许资源的时间共享，这增加了效率。

这些资源的价格在预订时引用，并根据（算法估计的）需求和供应变化动态计算。定时的更多用户灵活性可以通过降低价格自动获得回报。其基本思想是这些动态价格将调节需求，更好地利用云资源。这种机制还可以确保在高峰时期 - 需求根本无法满足 - 最“有价值”的工作就是那些可以运行而不是任意工作的工作。 ERA使用一种简单的投标模式，用户在每个工作请求中附加一个货币价值，指定他愿意为运行该工作而支付的最大金额。在这些高峰时段无法适应的工作所损失的价值量可作为购买额外云资源将获得的价值的量化，并且是云提供商采购决策流程的重要输入。

### 2.2 The Cloud Model （云模型）
The cloud in the ERA framework is modeled as an entity that sells multiple resources, bundled in configurations, and the capacity of these resources may change over time. The configurations, as well as the new concept of “virtual resources,” are designed to represent constraints that the cloud is facing, such as packing constraints. Specifically, the cloud is defined by: (1) a set of formal resources for sale (e.g., core or GB). We also allow for capturing additional constraints of the underlying infrastructure by using a notion of “virtual resources”; (2) a set of resource configurations: each configuration is defined by a bundle of formal resources (e.g., “ConfA” equals 4 cores and 7 GB),（These configurations are preset, but notice that ERA’s cloud model also supports the flexibility that each job can pick its own bundle of resources (as in YARN) by defining configurations of the basic formal resources (e.g., “ConfCore” equals a single core).） and is also associated with a bundle of actual resources that reflects the average amount the system needs in order to supply the configuration. The actual resources will typically be larger than the formal resources. The gap is supposed to model the overhead the cloud incurs when trying to allocate the formal amount of resources within a complex system. The actual resources can be composed of formal as well as virtual
resources; (3) inventory: the amount of resources (formal and virtual) over time; (4) time definitions of the cloud (e.g., the precision of time units that the cloud considers).

ERA框架中的云模拟为销售多种资源的实体，捆绑在配置中，这些资源的容量可能会随时间而改变。这些配置以及“虚拟资源”的新概念旨在表示云所面临的约束条件，如包装约束。具体而言，云定义为：（1）一组待售的正式资源（例如核心或GB）。我们还允许通过使用“虚拟资源”的概念来捕获底层基础设施的其他限制; （2）一组资源配置：每个配置由一组正式资源定义（例如，“ConfA”等于4个核心和7 GB），（这些配置是预设的，但注意到ERA的云模型也支持通过定义基本正式资源的配置（例如，“ConfCore”等于单个核心），每个工作都可以选择自己的资源束（如YARN中那样）。并且还与反映平均量的一批实际资源相关联系统需要为了提供配置。实际资源通常比正式资源要大。这个差距应该模拟云在尝试在复杂系统中分配正式的资源量时会产生的开销。实际资源可以由正式和虚拟组成资源; （3）库存：随着时间的推移资源量（正式和虚拟）; （4）云的时间定义（例如，云考虑的时间单位的精确度）。

### 2.3 The ERA Architecture （ERA架构）
The ERA system is designed as a smart layer that lies between the user and the cloud scheduler, as shown in Figure 1. The system receives a stream of job reservation requests for resources arriving online from users. Each request describes the resources the user wishes to reserve and the time frame in which these resources are desired, as well as an economic value specifying the maximal price the user is willing to pay for these resources. ERA grants a subset of these requests with the aim of maximizing total value (and/or revenue). The interface with these user requests is described in Section 2.3.1.

ERA interfaces with the cloud scheduler to make sure that the reservations that were granted actually get the resources they were promised. ERA instructs the cloud how it should allocate its resources to jobs, and the cloud should be able to follow ERA’s instructions and (optionally) to provide updates about its internal state (e.g., the capacity of available resources), allowing ERA to re-plan and optimize. ERA’s interface with the cloud scheduler is described in Section 2.3.2.

The architecture encapsulates the logic of the scheduling and pricing in the algorithm module. The algorithms use the prediction module to compute prices dynamically based on the anticipated demand and supply. This architecture gives the ability to change between algorithms and to apply different learning methods. ERA’s internal interface with the algorithmic components is described in Section 3.

ERA系统被设计为位于用户和云计划程序之间的智能层，如图1所示。系统从用户处接收到线上到达资源的作业预留请求流。每个请求都会描述用户希望保留的资源以及期望这些资源的时间范围，以及指定用户愿意为这些资源支付的最高价格的经济价值。 ERA授予这些请求的一部分，以实现总价值（和/或收入）最大化。第2.3.1节描述了这些用户请求的接口。

ERA与云调度程序进行接口，以确保授予的预留实际上获得了他们承诺的资源。 ERA指导云如何将资源分配给作业，并且云应该能够遵循ERA的指示和（可选地）提供关于其内部状态（例如可用资源的容量）的更新，从而允许ERA重新计划并进行优化。第2.3.2节描述了ERA与云调度程序的接口。

该体系结构在算法模块中封装了调度和定价的逻辑。算法使用预测模块根据预期的需求和供应动态计算价格。这种架构能够在算法之间进行切换，并应用不同的学习方法。第3节介绍了ERA与算法组件的内部接口。

#### 2.3.1 ERA-User Interface （ERA用户接口）
The ERA-User interface handles a stream of reservation requests that arrive online from users, and determines which request is accepted and at which price.

ERA用户界面处理从用户在线到达的预定请求流，并确定接受哪个请求以及以哪个价格。



*The Bidding Description Language BDL 投标说明语言*

Each reservation request bids for resources according to ERA’s bidding description language – an extension of the reservation definition language formally defined in [10]. The bid is composed of a list of resource requests and a maximum willingness to pay for the whole list. Each resource request specifies the configurations of resources that are requested, the length of time for which these are needed, and a time window [arrival, deadline). All the resources must be allocated after the arrival time (included) and before the deadline (excluded). For example, a resource request may ask for a bundle of 3 units of ConfA and 2 units of ConfB, for a duration of 2 hours, sometime between 6AM and 6PM today. Each configuration is composed of one or more resources, as described in Section 2.2.

By supporting a list of resource requests, ERA allows the description of more complex jobs, including the ability to break each request down to the basic resource units allowing for MapReduce kinds of jobs, or to specify order on the requests to some degree. The current ERA algorithms accept a job only if all of the resource requests in the list can be supplied; i.e., they apply the AND operator between resource requests. More sophisticated versions may allow more complex and rich bidding descriptions, e.g., support of other operators or recursive bids. For clarity of presentation, in this paper we present ERA in the simple case, where the reservation request is a single resource request, and there is only a single resource rather than configurations of multiple resources.

每个预留请求根据ERA的投标描述语言为资源投标 - 这是[10]中正式定义的预留定义语言的扩展。投标由资源请求列表和最大支付整个列表的意愿组成。每个资源请求指定请求的资源配置，需要这些资源的时间长度以及时间窗口[到达，截止日期）。所有资源必须在到达时间（包括）和截止日期之前（排除）分配。例如，资源请求可能会要求提供3个ConfA单元和2个ConfB单元，时长为2个小时，现在是上午6点到下午6点之间。每个配置由一个或多个资源组成，如2.2节所述。
通过支持资源请求列表，ERA允许描述更复杂的作业，包括将每个请求分解为允许MapReduce类型作业的基本资源单元的能力，或者在某种程度上指定请求的顺序。只有当列表中的所有资源请求都可以提供时，当前的ERA算法才接受作业;即它们在资源请求之间应用AND运算符。更复杂的版本可能允许更复杂和更丰富的投标说明，例如支持其他运营商或递归投标。为了表达清晰，在本文中，我们介绍简单情况下的ERA，其中预留请求是单个资源请求，并且只有一个资源而不是多个资源的配置。



*The makeReservation method （makeReservation方法）*

The interface with the user reservation requests is composed of the single “makeReservation” method, which handles a job reservation request that is sent by some user. Each reservation request can either be accepted and priced or declined by ERA. The basic input parameters to this method are the job’s bid and the identifier of the job. The bid encapsulates a list of resource requests along with the maximum price that the user is willing to pay in order to get this request (as described above). The output is an acceptance or rejection of the request, and the price that the user will be charged for fulfilling his request in case of acceptance.（Alternatively, the system may allow determining the payment after running is completed (depending on the system load at that time), or may allow flexible payments that take into account boththe amount of resources reserved and the amount of resources actually used.）

The main effect of accepting a job request is that the user is guaranteed to be given the desired amount of resources sometime within the desired time window. An accepted job must be ready to use all requested resources starting at the beginning of the requested window, and the request is considered fulfilled as long as ERA provides the requested resources within the time window.

与用户预订请求的接口由单个“makeReservation”方法组成，该方法处理由某个用户发送的作业预留请求。每个预订请求都可以被接受和定价或由ERA拒绝。该方法的基本输入参数是作业的出价和作业的标识符。该出价封装了资源请求列表以及用户愿意为获得此请求而支付的最高价格（如上所述）。输出是对请求的接受或拒绝，以及用户在接受时为履行他的请求而收取的价格（或者，系统可允许在运行完成后确定支付（取决于系统在那时候），或者可以允许灵活的支付，既考虑到保留的资源量又考虑实际使用的资源量。）

接受工作请求的主要作用是确保用户在预期的时间范围内某个时间获得所需的资源量。已接受的作业必须准备好在请求的窗口开始时使用所有请求的资源，并且只要ERA在时间窗口内提供请求的资源，就认为该请求已满足。

#### 2.3.2 ERA-Cloud Interface （ERA云端接口）

The interface between ERA and the cloud-scheduler is composed of two main methods that allow the cloud to get information about the allocation of resources it should apply at the current time, and to provide ERA with feedback regarding the actual execution of jobs and changes in the cloud resources.

ERA和云调度程序之间的接口由两种主要方法组成，它们允许云获得当前应该应用的资源分配信息，并向ERA提供有关作业实际执行情况和云资源。



*The getCurrentAllocation method*

This is the main interface with the actual cloud scheduler. The cloud should repeatedly call this method (quite often, say, every few seconds) and ask ERA for the current allocations to be made.（For performance, it is also possible to replace this query with an event-driven scheme in which ERA pushes an event to the cloud scheduler when the allocations change.）The method returns an allocation, which is the list of jobs that should be instantaneously allocated resources and the resources that should be allocated to them. In the simple case of a single resource, it is a list of “job J should now be getting W resources.” The actual cloud infrastructure should update the resources that it currently allocates to all jobs to fit the results of the current allocation returned by this query. This new allocation remains in effect until a future query returns a different allocation. It is the responsibility of the underlying cloud scheduling system to query ERA often enough, and to put these new allocations into effect ASAP, so that any changes are effected with reasonably small delay. The main responsibility of the ERA system is to ensure that the sequence of answers to this query reflects a plan that can accommodate all accepted reservation requests.

The main architectural aspect of this query is to make the interface between ERA and the cloud system narrow, such that it completely hides the plan ERA has for future allocation. It is assumed that the cloud has no information on the total requirements of the jobs, and follows ERA as accurately as possible.

*getCurrentAllocation方法*

这是实际云调度程序的主界面。云应该反复调用这个方法（通常是说，每隔几秒），并要求ERA进行当前的分配（对于性能，也可以用ERA推动的事件驱动方案替换该查询当分配发生更改时，该事件将发送到云计划程序）。该方法返回一个分配，该分配是应当立即分配资源的作业的列表以及应分配给它们的资源。在单一资源的简单情况下，它是“作业J现在应获取W资源”的列表。实际云基础架构应更新当前分配给所有作业的资源，以适应当前分配结果这个查询。这个新分配保持有效，直到将来的查询返回不同的分配。底层云调度系统有责任经常查询ERA，并尽快将这些新的分配生效，以便任何更改都以合理的小延迟进行。 ERA系统的主要职责是确保该查询的答案顺序反映了可以适应所有接受的预订请求的计划。

此查询的主要体系结构方面是使ERA和云系统之间的接口变窄，从而完全隐藏了ERA用于未来分配的计划。假定云没有关于工作总需求的信息，并且尽可能准确地遵循ERA。



*The update method (optional usage)*

The cloud may use this method to periodically update ERA with its actual state. Using this method is important since the way resources are actually used in real time may be different from what was planned for. For example, some processors may fail or be taken offline. Most importantly, it is expected that most jobs will use significantly less resources than what they reserved (since by trying to ensure that they have enough resources to usually complete execution, they will probably reserve more than they actually use). The ERA system should take this into account and perhaps re-plan.

The simple version of the cloud feedback includes: (1) changes in the current resources under the cloud’s management (e.g., if some computers crashed); (2) the current resource consumption; (3) termination of jobs; (4) the number of waiting processes of each job, which specifies how many resources the job could use at this moment, if the job were allocated an infinite amount.

*更新方法（可选用法）*

云可以使用此方法定期更新其实际状态的ERA。使用这种方法很重要，因为实际使用资源的方式可能与计划中的不同。例如，某些处理器可能会失败或处于脱机状态。最重要的是，预计大多数工作将使用比他们保留的资源少得多的资源（因为试图确保他们有足够的资源来完成执行，他们可能会预留比实际使用更多的资源）。电子逆向拍卖系统应考虑到这一点，并可能重新计划。

云反馈的简单版本包括：（1）云管理下当前资源的变化（例如，如果某些计算机崩溃）; （2）当前的资源消耗; （3）终止工作; （4）每个作业的等待进程的数量，如果该作业被分配了无限量，则指定该作业此时可以使用多少资源。


## 3. ALGORITHMS 算法

The internal algorithmic implementation of ERA is encapsulated in separate components – the algorithm and the prediction components – in a flexible “plug and play” design, allowing to easily change between different implementations to fit different conditions and system requirements. The algorithm component is where the actual scheduling and pricing of job requests are performed. The algorithm may use the prediction component in order to get the market prices or the estimated demand, and the ERA system updates the prediction component online with every new request.

ERA的内部算法实现采用灵活的“即插即用”设计封装在单独的组件（算法和预测组件）中，允许在不同的实现之间轻松切换以适应不同的条件和系统要求。算法组件是执行作业请求的实际调度和定价的地方。该算法可以使用预测分量来获得市场价格或估计需求，并且ERA系统使用每个新请求在线更新预测分量。

### 3.1 Scheduling and Pricing Algorithms Interface 调度和定价算法接口

The ERA algorithm is an online scheduling and pricing algorithm that provides the logic of an ERA system. The ERA system forwards queries arriving from users and from the cloud to be answered by the algorithm, and so the internal interface between ERA and the algorithm is similar to the external ERA interface (described in Sections 2.3.1 and 2.3.2), except that it abstracts away all the complexities of interfacing with the external system. The main change between these two interfaces is that the algorithm is not given the bids (the monetary value) of the reservation requests, and must decide on the price independently of the bid. It can only make a one-time comparison against the value, and the request is accepted as long as the value is not smaller than the price. Thus, the architecture enforces that the algorithm will be monotonic in value (as it sets a threshold price for winning), creating an incentive compatible mechanism with respect to the value; i.e., the resulting mechanism is truthful by design.

The scheduling and pricing of a new job is performed in the makeReservation method. As described in detail in Section 2.3.1, the input to this method is a reservation request of the form “I need W cores for T time units, somewhere in the time range [Arrival,Deadline), and will pay at most V for it.” The answers are of the form “accept/reject” and a price P in case of acceptance. The algorithm should also keep track of its planned allocations to actually tell the cloud infrastructure when to run the accepted jobs upon a getCurrentAllocation query, and re-plan and optimize upon an update query (see Section 2.3.2).

ERA算法是提供ERA系统逻辑的在线调度和定价算法。 ERA系统转发来自用户和云端的查询，并由算法应答，因此ERA和算法之间的内部接口类似于外部ERA接口（在2.3.1和2.3.2中描述），除它将所有与外部系统接口的复杂性抽象出来。这两个接口之间的主要变化是算法没有给出预订请求的出价（货币值），并且必须独立于出价来决定价格。它只能对价值进行一次性比较，只要价值不低于价格，就会接受请求。因此，该体系结构强制该算法在价值上是单调的（因为它为获胜设定了阈值价格），从而创建了与价值相关的激励兼容机制;即由此产生的机制在设计上是真实的。

新工作的计划和定价在makeReservation方法中执行。如2.3.1部分详细描述的那样，该方法的输入是“我需要T时间单位的W核心，在某个时间范围[到达，截止日期]的预订请求，并且最多支付V它“。答案的形式是”接受/拒绝“，如果接受，答案是价格P.该算法还应跟踪其计划分配，以实际告诉云基础设施何时在getCurrentAllocation查询时运行已接受的作业，并在更新查询时重新计划和优化（请参阅第2.3.2节）。

*Basic Econ Scheduling 基本经济调度*

The Basic Econ Scheduling (Algorithm 1) is our basic implementation of an ERA algorithm. Whenever a new job request arrives, the algorithm dynamically sets a price for each time and each unit of the resource (e.g., second*core), and the total price is the sum over these unit prices for the requested resources.（In case of multiple resources, the simple generalization is to set the total price additively over the different types of resources. We choose to focus on additive pricing due to its simplicity and good economic properties (e.g., splitting a request is never beneficial).）It then schedules the job to start at the cheapest time within its requested window that fits the request, as long as the job’s value is not lower than the computed total price. To determine the price of a resource in a specific time unit t in the future, the algorithm takes into account the amount of resources already promised as well as its prediction for future demand for that time unit. Essentially, the price is set to reflect the externalities imposed on future requests due to accepting this job, according to the predicted demand. The prediction of demand is encapsulated in the prediction component we will discuss in the next section.

Note that this simple algorithm gives up the flexibility to preempt jobs (swap jobs in and out) and instead allocates to each job a continuous interval of time with a fixed starting time. It also allocates exactly the W requested cores concurrently instead of trading off more time for less parallel cores. We chose to give up these flexibilities in the basic implementation, although they are supported by the ERA API, in order to isolate concerns: this choice separates the algorithmic issues (which are attacked only in a basic way) from pricing issues (which are dealt with) and from learning issues. In addition, such schedules are robust and applicable under various real-world constraints, and in other cases they may simply be suboptimal and serve as benchmarks.

基本经济调度（算法1）是ERA算法的基本实现。每当新的工作请求到达时，算法会动态地为每个时间和资源的每个单位（例如second * core）设置一个价格，总价格就是所请求资源的这些单位价格的总和（如果是多种资源，简单的概括是将总价格累加地设置在不同类型的资源上。由于其简单性和良好的经济属性，我们选择关注附加定价（例如，拆分请求永远不会有好处））。然后，它会安排工作在其请求的窗口内的最便宜的时间开始，以符合请求，只要该工作的价值不低于计算出的总价格。为了确定将来特定时间单元t中资源的价格，该算法考虑已经承诺的资源量以及对该时间单元未来需求的预测。本质上，根据预测的需求，价格被设定为反映由于接受这份工作而对未来请求施加的外部性。需求预测被封装在我们将在下一节讨论的预测组件中。

注意，这个简单的算法放弃了抢占作业的灵活性（交换作业进出），而是以固定的开始时间为每个作业分配连续的时间间隔。它还同时分配W请求的内核，而不是为更少的并行内核换取更多时间。我们选择放弃基本实施中的这些灵活性，尽管它们受ERA API支持，以便隔离关注点：该选择将算法问题（仅以基本方式进行攻击）与定价问题（将其分别处理与）和学习问题。另外，这样的时间表在各种现实世界约束条件下是稳健的和可应用的，并且在其他情况下，它们可能仅仅是不理想的并且用作基准。

![](http://opkk27k9n.bkt.clouddn.com/18-3-3/60803933.jpg)

### 3.2 Demand Prediction Interface 需求预测接口

The prediction component is responsible for providing an estimation of demand at a future time t at any given price, given the current time. Since the inverse function is what we really need, our actual interface provides that inverse function:（Yet, we present the predictors using both the demand function and its inverse. Moving between the two is straightforward.） given a future time t, the current time, and a quantity of demand q, it returns the highest price such that the demand that arrives from the current time till t, at this price, is equal to the specified quantity q.

In general, one cannot expect future demand to be determined deterministically – thus a prediction would, in its most general form, be required to specify a probability distribution over prices that will result in selling the specified quantity. As such an object would be hard to work with, our basic implementation simplifies the prediction problem, and requires the predictor to only specify a single price for each demand quantity, as if demand is deterministic. Such an approach is justified when the total demand is a result
of the aggregation of a large number of independent requests. In that case the demand will be concentrated and the single expected price will reasonably approximate the price distribution.

在给定当前时间的情况下，预测组件负责提供在任意给定价格的未来时间t的需求估计。由于反函数是我们真正需要的，我们的实际界面提供了反函数:(然而，我们使用需求函数及其逆函数来呈现预测器。在两者之间移动是直截了当的）给定未来时间t，当前时间和需求量q，它返回最高价格，使得从当前时间到达t的需求在此价格等于指定数量q。

一般来说，不能指望确定性地确定未来的需求 - 因此，预测将以其最一般的形式被要求指定将导致出售指定数量的价格概率分布。由于这样的对象很难合作，所以我们的基本实现简化了预测问题，并且要求预测器仅为每个需求量指定单一价格，就好像需求是确定性的。当总需求是结果时，这种方法是合理的
的大量独立请求的汇总。在这种情况下，需求将集中并且单个预期价格将合理地接近价格分布。

*Data-based predictors: prediction based on historic data 基于数据的预测变量：基于历史数据的预测*

ERA’s predictor – the demand oracle – builds its estimations based on historic data of job requests. It gets as input a list of past requests, and learns, for every time t, the demand curves (i.e., the demand as a function of price) according to the demand in the list. Of course, this approach presents multiple challenges: first, there is the “cold start” problem – as ERA defines a new interface for job requests, there are no past requests of the form that ERA can use to learn. Second, the success of the prediction depends on the ability to determine cycles in the demand, such as day-night or days of the week. In addition, the learning methods must also overcome sampling errors and address the non-deterministic nature of the demand (as discussed above).

Our first implementation of a data-based predictor puts these challenges aside and aims to suggest an approach to address an additional major challenge: the time flexibility of jobs. Essentially, the problem is that we expect the predictor to provide the instantaneous demand, while in ERA the requests are for resources for some duration, within a window of time that is usually longer than the duration. Thus, we should answer the following question: how should a job request affect the predicted demand in its requested time window?

The naive approach would be to spread the demand over the window, e.g., a request of 10 cores for 5 minutes over a window of 50 minutes would contribute 1 core demand in each of the 50 minute window. However this may not reflect the actual demand we should expect. For example, consider the input of low-, medium-, and high-value jobs. Each type asks for 100% of the capacity, where the high-value jobs can run only during the day and the low- and the medium-value jobs can run either day or night. Using the spreading approach we count the demand of the high-value jobs at day, and spread the low- and medium-value jobs over day and night, such that at night we obtain a demand of 50% of the low- and 50% of the medium-value jobs. Using this demand gives the impression that we can fill only half of the capacity using the medium-value jobs at night, and so we will set the price to be too low, and will accept low-value jobs at the cost of declining medium-value ones.

We suggest that this problem can be overcome by taking a different approach based on the LP relaxation of the problem. The LP-based predictor runs a linear program, offline, to find the optimal (value-maximizing) fractional allocation over past requests, and predicts the demand at time t using the fractional optimal allocation at that time. Note that this LP requires many variables – one variable for every job and every time in the job’s time window, and the number of degrees of freedom may be large, and so one may suspect that the predicted demand will be very different at different times that are experiencing essentially the same demand. Our preliminary empirical tests suggest that this LP-based approach is stable, yet future work should test this further and establish theoretical justifications for the approach.

RA的预测器 - 需求预测器 - 根据工作请求的历史数据构建其估计。它将过去请求列表作为输入，并且在每个时间t根据列表中的需求获得需求曲线（即需求作为价格的函数）。当然，这种方法带来了多种挑战：首先，存在“冷启动”问题 - 由于ERA为工作请求定义了一个新的界面，因此ERA可以用来学习的形式没有过去的要求。其次，预测的成功取决于确定需求周期的能力，例如一天中的某一天或一周中的几天。此外，学习方法还必须克服抽样误差并解决需求的非确定性性质（如上所述）。

我们首次实施的基于数据的预测工具将这些挑战放在一边，旨在提出一种解决额外重大挑战的方法：工作的时间灵活性。本质上，问题在于我们期望预测器提供即时需求，而在ERA中，请求是在一段时间内（通常比持续时间更长）的时间内的资源。因此，我们应该回答以下问题：求职申请如何影响其所要求的时间窗内的预测需求？

天真的做法是将需求扩散到窗口上，例如，在50分钟的窗口上10分钟内请求5分钟将在50分钟窗口的每一个窗口中贡献1个核心需求。但这可能并不能反映我们应该预期的实际需求。例如，考虑低价值，中等价值和高价值工作的投入。每种类型都要求100％的能力，高价值工作只能在白天运行，低价值工作和中等价值工作可以在白天或晚上运行。我们使用扩展方法计算白天高价值工作的需求，并在白天和黑夜分散低价值和中等价值的工作，以便在晚上我们获得50％的低价和50％的中等价值工作。使用这种需求给人的印象是，我们只能在夜间使用中等价值的工作来填补一半的产能，因此我们将把价格定得太低，并且会接受低价值的工作，价值的。

我们建议这个问题可以通过采用基于问题的LP松弛的不同方法来克服。基于LP的预测器离线运行线性程序，以查找过去请求的最优（最大化值）分数分配，并使用当时的分数最优分配预测时间t的需求。请注意，这个LP需要许多变量 - 每个作业和每个作业的时间窗口都有一个变量，自由度的数量可能很大，因此人们可能会怀疑预测的需求在不同的时间会有很大的不同正在经历基本相同的需求。我们的初步实证测试表明，这种基于LP的方法是稳定的，但未来的工作应该进一步测试，并为该方法建立理论上的理由。

## 4. THE ERA SYSTEM AND SIMULATIONS （ERA系统和仿真）

ERA is a complete working system: it is implemented as a software package that provides the interfaces described above together with basic implementations of the pricing, scheduling, and prediction algorithms, which are pluggable and can be extended or replaced. In addition, the system contains a simulation platform that can simulate the execution of an algorithm given a sequence of job requests and a model of the underlying cloud, using exactly the same core code that is interfaced with the real cloud and users. See the system architecture in Figure 1 and a screen-shot of the simulator in Figure 2.

We have performed extensive runs of ERA within the simulator as well as proof-of-concept runs with two cloud systems: Hadoop/YARN and Microsoft’s Azure Batch simulator. Next we present a few of these runs to demonstrate the large potential gains of moving from the simple cloud-pricing systems, like the ones currently in use, to ERA – the Economic Resource Allocation system – and to demonstrate the ability of the ERA system to integrate with existing
cloud systems.

ERA是一个完整的工作系统：它被实现为一个软件包，提供上述接口，以及可插入，可扩展或替换的定价，调度和预测算法的基本实现。此外，该系统还包含一个仿真平台，可以使用与真实云和用户接口相同的核心代码，模拟执行一系列作业请求和底层云模型的算法。请参阅图1中的系统架构和图2中模拟器的屏幕截图。

我们在模拟器中执行了广泛的ERA运行，以及两个云系统的概念验证运行：Hadoop / YARN和Microsoft的Azure Batch模拟器。接下来，我们将展示其中的一些展示，以展示从简单的云定价系统（如目前正在使用的系统）向ERA（经济资源分配系统）转变的巨大潜在收益，并展示ERA系统的能力与现有的整合
云系统。


*The importance of economic allocation 经济分配的重要性*

We first demonstrate the ability of the ERA system to improve the efficiency of cloud resource usage significantly, by considering the jobs’ values. We use the simulator with input of jobs that were sampled according to distributions describing a large-scale MapReduce production trace at Yahoo [9], after some needed modifications of adding deadlines and values that were not included in the original trace. In this input, there are 6 classes of jobs, and about 1,400–1,500 jobs of each class. Jobs of class “yahoo-5” have the largest average size, and we set them to have a low average value per unit of ` $1`, while we set jobs of all other classes to have a high value per unit of $10, to model high-value production jobs. The cluster is way too small to fit all jobs.

We compare ERA’s Basic-Econ scheduling algorithm with a greedy algorithm that does not take into account the values of the jobs, but instead charges a fixed price per resource unit, and that schedules the job to run at the earliest possible time within its requested time window. The simulation shows that the greedy algorithm populates most of the cluster with the large, low-value jobs (of class yahoo-5) and results in a low efficiency of only 10% of the total requested value. In sharp contrast, ERA’s Basic-Econ algorithm, which is aware of the values of the jobs and uses dynamic pricing to accept the higher value jobs, achieves 51% of the requested value (note that getting 100% is not possible as the cloud is too small to fit all jobs).

我们首先通过考虑作业的价值来证明ERA系统显着提高云资源使用效率的能力。我们使用模拟器输入根据描述Yahoo [9]上的大规模MapReduce生产追踪的分布进行抽样的作业，并在添加了未包含在原始追踪中的期限和值后进行一些必要的修改。在这个投入中，有6类职位，每类职位约有1,400-1,500个职位。 “yahoo-5”级别的职位的平均规模最大，我们将其设置为每单位1美元的平均价值较低，而我们将所有其他职位的职位设置为每单位10美元的高价值，价值生产工作。该集群太小，不适合所有的工作。

我们将ERA的Basic-Econ调度算法与贪婪算法进行比较，该算法没有考虑作业的值，而是对每个资源单位收取固定价格，并安排作业在要求的时间内尽早运行窗口。仿真表明，贪婪算法使用大量低价值作业（类yahoo-5）填充大部分集群，并且导致效率低于总请求值的10％。与此形成鲜明对比的是，ERA的Basic-Econ算法能够了解作业的价值，并使用动态定价来接受更高价值的作业，达到了所请求价值的51％（请注意，由于云计算是100％是不可能的太小而不适合所有工作）。

*ERA–Rayon integration ERA-Rayon集成*

We next demonstrate that it is feasible to integrate ERA with a real cloud system by showing that the cloud succeeds in running real jobs using ERA. In addition, we show that the ERA simulator provides a good approximation to the outcome of the real execution.

We have fully integrated ERA with Rayon [10], which is a cloud system that handles reservations for computational resources, and is part of Hadoop/YARN [31] (aka MapReduce 2.0). The integration required, first, that we introduce economic considerations into the Rayon system, as Rayon’s original reservation mechanism did not consider the reservations’ monetary valuations. Next, we plugged ERA’s core code into Rayon’s reservation and scheduling process, by adding a layer of simple adapter classes that bridge between ERA’s and Rayon’s APIs. The bridging layer configured Rayon to completely follow ERA’s instructions via the getCurrentAllocation method (see Section 2.3.2), but made one extension to this query: it added an “empty allocation” (i.e., allocation of zero resources), for jobs that are during their reservation time-window but are currently not allocated resources. Rayon opened a queue for each job that was returned by ERA, including jobs with an empty allocation, and thus it was able to run jobs earlier than they were scheduled when it was possible.

We tested the integration by using a workload of MapReduce jobs that we generated using the Gridmix(http://hadoop.apache.org/docs/r1.2.1/gridmix.html) platform. The jobs read and wrote synthetic data from files of 100 GB created for this purpose. Eight hundred and fifteen jobs were processed, all of which finished successfully. They arrived during a period of one hour, asked on average for 3 GB memory, for a duration of 60 seconds on average (σ = 6 seconds). The cluster consisted of 3 nodes, of 80 GB memory each. Rayon’s resource manager was configured to use ERA with the simplest greedy algorithm (described above) that allocates a single resource – GB of memory (as the version of Rayon at the time allocated only memory).

We ran the same job workload in the ERA simulator, with the same greedy algorithm, and a cloud model that communicates with ERA every second with no failures. The comparison between these two runs – over Rayon (Hadoop) system and in the simulator – shows that the simulator gives a good approximation to the performance of ERA on a cloud system. We found that jobs were scheduled and running on approximately similar points in time and had similar durations. The main difference between the two runs is that while the simulator assigns jobs a constant capacity throughout their (simulated) execution, the real cluster changes their capacity according to various system considerations that are out of ERA’s control. The total allocation obtained in these two runs (GB*sec) was similar: 76,730 using the simulator vs. 77,056 in the real cloud.


我们接下来证明，通过显示云成功地使用ERA运行实际作业，将ERA与真正的云系统集成是可行的。另外，我们表明ERA模拟器提供了一个很好的逼近真实执行的结果。

我们将ERA与Rayon [10]完全集成，这是一个处理计算资源保留的云系统，是Hadoop / YARN [31]（又名MapReduce 2.0）的一部分。整合所需要的首先是我们将经济考虑引入人造丝系统，因为人造丝的原始保留机制没有考虑保留的货币估值。接下来，我们通过添加一层简单的适配器类，将ERA的核心代码插入到Rayon的预留和调度过程中，这些适配器类在ERA和人造丝的API之间建立了桥梁。桥接层通过getCurrentAllocation方法（参见第2.3.2节）配置了Rayon以完全遵循ERA的指令，但对此查询做了一个扩展：它为作业添加了“空分配”（即分配零资源）在他们的预订时间窗口期间，但目前没有分配资源。 Rayon为ERA返回的每项工作打开了一个队列，其中包括分配空置的工作，因此它能够在可能的时候提前执行工作。

我们通过使用我们使用Gridmix(http://hadoop.apache.org/docs/r1.2.1/gridmix.html)平台生成的MapReduce作业的工作负载来测试集成。这些作业从为此创建的100 GB文件中读取和写入合成数据。共处理了八百一十五份工作，所有这些工作都顺利完成。他们在一小时内到达，平均要求3 GB的内存，平均持续60秒（σ= 6秒）。该群集由3个节点组成，每个节点80 GB。 Rayon的资源管理器被配置为使用ERA和最简单的贪婪算法（如上所述）来分配单个资源 - GB内存（作为当时分配给内存的Rayon版本）。

我们使用相同的贪婪算法在ERA模拟器中运行相同的作业工作负载，并且每秒都会与ERA进行通信并且没有发生故障的云模型。这两次运行之间的比较（在Rayon（Hadoop）系统和模拟器中）显示模拟器可以很好地近似于云系统上的ERA的性能。我们发现，工作安排在大致类似的时间点上，并且具有相似的持续时间。两次运行之间的主要区别在于，尽管模拟器在整个（模拟）执行过程中为作业分配了一个恒定容量，但真正的集群根据ERA控制范围之外的各种系统考虑因素改变其容量。在这两次运行中获得的总分配（GB \*秒）类似：使用模拟器的总数为76,730，而真实云数量为77,056。



*Testing Azure Batch 测试Azure Batch*

The next set of simulations shows the advantage of using ERA over existing algorithms when applied on a cloud scale. In a typical cloud environment, we cannot expect one instance of ERA to have complete control of millions of cores. Thus, our goal here is to evaluate whether ERA will work with a subset of cores in a region, even while the underlying resource availability is constantly changing.

The simulations were of a datacenter consisting of 150K cores. ERA was given access to 20% of the resources and the remaining 80% were allocated to non-ERA requests, which were modeled using the standard Azure jobs. This means that resources were constantly being allocated/freed in the underlying region and ERA had to account for this. The 20% of the resources under ERA’s management came from the pre-emptible resources, but the design does not restrict its use to pre-emptible resources alone. ERA itself was run as a layer on top of the Azure Batch simulator, which simulates batch workloads on top of the Azure simulator of Microsoft.

ERA’s Basic-Econ scheduling algorithm was experimented relative to two other algorithms: (1) the on-demand algorithm, which accepts jobs if there are enough available resources to start and run them (availability is checked only for the immediate time, ignoring the duration that the resources are requested). It schedules accepted jobs to run immediately and charges a fixed price; (2) the greedy (“FirstFit”) algorithm (described above), which charges the fixed, discounted, price of 65% of the non-pre-emptible resources price.

A common practice in the industry is to bound the maximal discount over non-pre-emptible machines. Accordingly, in our experiments ERA’s Basic-Econ algorithm was restricted so that the price would be no higher than the non-pre-emptible jobs and would give no more than 35% discount. Several variants of the econ algorithm were explored: (1) using either a linear predictor that is based on prior knowledge of the job distributions, or a predictor that uses past observations; (2) with or without an exponential penalty for later scheduling. Each of the variants was tested at a different capacity of the algorithm’s use, so that the higher the capacity the fewer the resources that remained as spares for re-running failed jobs.

All jobs in the simulation workloads requested a time-window that started at their request-time (i.e., jobs did not reserve in advance). As ERA was getting 20% of the resources, we wanted to evaluate two measure metrics: (1) late-job percentage: this is the percentage of jobs that finished later than their deadlines; (2) accepted revenue: as we can charge only for jobs that are accepted, the better the algorithm, the more jobs we can accept. Figure 3 shows that ERA’s econ algorithm dominates the other algorithms in terms of these two desired measures.


下一组仿真显示了当应用于云规模时使用ERA优于现有算法的优势。在典型的云环境中，我们不能期望有一个ERA实例能够完全控制数百万个核心。因此，我们的目标是评估ERA是否能够与某个地区的核心子集一起工作，即使在底层资源可用性不断变化的情况下。

仿真是由150K内核组成的数据中心。 ERA获得20％的资源，其余80％分配给非ERA请求，这些请求是使用标准Azure作业建模的。这意味着资源在基础地区不断分配/释放，ERA必须对此进行解释。 ERA管理的资源中有20％来自可预先安排的资源，但该设计并未将其用途仅限于可预先安排的资源。 ERA本身作为Azure批处理模拟器的顶层运行，该模拟器模拟Microsoft Azure模拟器上的批处理工作负载。

ERA的Basic-Econ调度算法相对于其他两种算法进行了实验：（1）按需算法，如果有足够的可用资源来启动和运行它们，则接受作业（可用性仅在立即时间检查，忽略持续时间资源被请求）。它安排接受的工作立即运行并收取固定价格; （2）贪婪（“FirstFit”）算法（如上所述），该算法收取非固定资源价格的固定折扣价格的65％。

业内普遍的做法是限制非优先机器的最大折扣。因此，在我们的实验中，ERA的Basic-Econ算法受到限制，因此价格不会高于非抢先式工作，并且折扣不会超过35％。探索了econ算法的几种变体：（1）使用基于工作分布的先验知识的线性预测器或使用过去观测的预测器; （2）对于以后的调度有或没有指数惩罚。每种变体都在算法使用的不同容量下进行测试，以便容量越高，作为重新运行失败作业的备件所剩的资源越少。

仿真工作负载中的所有作业都请求一个从其请求时间开始的时间窗口（即作业没有预先保留）。随着ERA获得20％的资源，我们希望评估两个度量指标：（1）晚工作比例：这是完成时间晚于截止日期的工作的百分比; （2）接受收入：因为我们只能为被接受的工作收取费用，算法越好，我们可以接受的工作就越多。图3表明ERA的econ算法在这两个期望的度量方面占优势。

![](http://opkk27k9n.bkt.clouddn.com/18-3-3/4429590.jpg)
Figure 3: ERA over Azure Batch – simulation results (axis scales removed). ERA’s econ algorithm dominates on demand and first-fit algorithms in terms of the two desired measures of revenue and percentage of late jobs.

图3：Azure Batch上的ERA - 模拟结果（删除轴比例）。 ERA的econ算法在需求和首次拟合算法方面占据主导地位，这主要取决于两种期望的收入指标和晚期工作比例。

## 5. GRAND CHALLENGES 巨大的挑战

Clearly, the main challenge is to get the ERA system integrated in a real cloud system, and interface with real paying costumers. Short of this grand challenge, there are many research challenges. In this section we describe several challenges of a practical and theoretical nature related to the ERA (Economic Resource Allocation) project.

显然，主要的挑战是将ERA系统集成到真正的云系统中，并与真正的付费用户进行交互。面临这一重大挑战，有许多研究挑战。在本节中，我们将介绍与ERA（经济资源分配）项目相关的实际和理论性质的几个挑战。

### 5.1 Job Scheduling 作业调度
There is a vast literature on job scheduling both in the stochastic and adversarial models. The most obvious related model is job scheduling with laxity, which is the difference between the arrival time and the latest time in which the job can be scheduled and still meet the deadline. The current issues that are raised by our framework give rise to new challenges in both domains. In our setting it is very reasonable to assume that any job requires only a small fraction of the total resources, and that the laxity is fairly large compared to the job size. An interesting realistic challenge is to have a job give a tradeoff between time (to run) and resources (number of machines), which depends on the degree of parallelism of the job. Another interesting challenge is to exhibit a model that interpolates between the stochastic model, which gives a complete model of the job arrival process, and the adversarial model, which does not make any assumptions. It would be nice to have a model that would require only a few parameters and be able to capture many arrival sequences. Finally, jobs of a reoccurring nature would be very interesting to study both in the stochastic and adversarial models.

在随机模型和敌对模型中都有大量关于作业调度的文献。最明显的相关模型是具有松弛性的工作时间安排，即到达时间与可安排工作的最近时间之间的差异，并且仍然符合截止日期。我们的框架提出的当前问题在两个领域都会带来新的挑战。在我们的情况下，假定任何工作只需要总资源的一小部分是非常合理的，而松散度与工作规模相比相当大。一个有趣的现实挑战是让工作在时间（运行）和资源（机器数量）之间进行权衡，这取决于工作的并行度。另一个有趣的挑战是展示一个模型，该模型在随机模型（它提供了一个完整的工作到达过程模型）和不做任何假设的敌对模型之间进行插值。如果有一个模型只需要几个参数并且能够捕获多个到达序列，那就太好了。最后，在随机模型和敌对模型中研究重复性质的工作将非常有趣。


### 5.2 Pricing 定价

In our model we assume that the user has both a clear deadline in mind and an explicit bound on the length of the job. It would be interesting to give a more flexible guarantee, which would help the user to set his preferences in a less conservative way. For example, one could allow the job to run after it exhausts its resources at a certain cost and at a slightly lower priority for a certain additional amount of time. Another similar guarantee is that the user would have his “preferred deadline” and his “latest deadline” with a guarantee that most jobs finish at the preferred deadline. All this is aimed at a more flexible Quality of Service (QoS) guarantee by the system. Pricing such complex guarantees is a significant practical and theoretical challenge.

From the theoretical side, it would be nice to give theoretical guarantees to our system. First, to show that the users have an incentive to report their information truthfully, and not to try and game the system, or at least achieve this approximately. Second, to show that the system reaches a satisfiable steady state (e.g., showing an appropriate equilibrium notion and a related price of anarchy).

在我们的模型中，我们假设用户既有明确的最后期限，也有明确的工作时间限制。提供更灵活的保证将会很有趣，这将有助于用户以较不保守的方式设定他的偏好。例如，一个人可以在一定的时间内以一定的成本和略低的优先级耗尽资源后开始工作。另一个类似的保证是，用户可以有他的“首选期限”和“最后期限”，并保证大多数工作在最佳期限内完成。所有这些都旨在为系统提供更加灵活的服务质量（QoS）保证。定价这种复杂的保证是一个重大的实践和理论挑战。

从理论上讲，为我们的系统提供理论保证会很好。首先，要表明用户有动机真实地报告他们的信息，而不是尝试和游戏系统，或者至少达到这一点。其次，为了表明系统达到可以满足的稳定状态（例如，显示适当的均衡概念和相关的无政府状态价格）。

### 5.3 Learning 学习

Our proposed framework requires a significant component of learning. Much of the learning depends on the observed time series from the past that would be used to predict future requests. A clear challenge in our setting is to accommodate seasonality effects (daily, such as day versus night; weekly, such as work week versus weekend; annual, such as holidays). Such challenges are well known in the time-series literature. A more interesting effect is that we have a system where the available resources and the demand are constantly growing, and the challenge is to bundle the two forecasts or somewhat separate them. It seems that our prediction model would need a more refined prediction than only the expected value, but for many of our forecasting applications we need to get more detailed information.

An additional uncertainty is that our system might be unable to see certain requests since the user decides that they were unlikely to be accepted and therefore never submitted them. For example, if a more important job is already rejected due to a low value, lessimportant jobs might be not submitted, and thus the prediction of the demand is even more challenging, given this partial information.

Finally, learning should not be limited only to the forecast of demand, but should also forecast the accuracy of the requests. Since in the current system we require that the job will not exceed its maximum length, it is likely to be a conservative estimate, and learning what is the “actual” demand might free significant resources.

我们提出的框架需要学习的重要组成部分。大部分学习依赖于过去观察到的时间序列，用于预测未来的请求。在我们的设置中，一个明显的挑战是适应季节性影响（每天，如白天与夜晚;每周，如工作周与周末;年度，如假期）。这些挑战在时间序列文献中是众所周知的。一个更有趣的结果是我们有一个可用资源和需求不断增长的系统，而挑战在于捆绑两个预测或者将它们分开。似乎我们的预测模型需要比预期值更精确的预测，但对于我们的许多预测应用程序，我们需要获得更详细的信息。

另一个不确定性是我们的系统可能无法看到某些请求，因为用户决定他们不太可能被接受，因此从未提交过。例如，如果一项较重要的工作由于价值较低而已被拒绝，那么较不重要的工作可能不会提交，因此，鉴于这种部分信息，对需求的预测更具挑战性。

最后，学习不应仅限于需求预测，还应该预测请求的准确性。由于在目前的系统中，我们要求工作不会超过其最大长度，所以这可能是一个保守的估计，而了解什么是“实际”需求可能会释放大量资源。

## 5.4 Robustness 鲁棒性
For any practical system to run it needs a significant level of robustness. Robustness should take into account both planned and unexpected failures in the various resources. Modeling this might be done as part of the greater challenge of a QoS guarantee. We should study what kind of an extreme-case guarantee can we give.

对于任何实际的系统来说，它都需要一个强大的稳健性水平。 健壮性应该考虑到各种资源中的计划和意外故障。 建模可以作为QoS保证更大挑战的一部分来完成。 我们应该研究一下我们可以提供什么样的极端情况保证。

## References 引用
[1] Apache Hadoop Project. http://hadoop.apache.org/.
[2] V. Abhishek, I. A. Kash, and P. Key. Fixed and market pricing for
cloud services. arXiv preprint arXiv:1201.5621, 2012.
[3] O. Agmon Ben-Yehuda, M. Ben-Yehuda, A. Schuster, and D. Tsafrir.
Deconstructing amazon ec2 spot instance pricing. ACM Transactions
on Economics and Computation, 1(3):16, 2013.
[4] Amazon. Amazon elastic mapreduce. At http://aws.amazon.
com/elasticmapreduce/.
[5] M. Armbrust, A. Fox, R. Griffith, A. D. Joseph, R. Katz, A. Konwinski,
G. Lee, D. Patterson, A. Rabkin, I. Stoica, et al. A view of cloud
computing. CACM, 53(4):50–58, 2010.
[6] Y. Azar, I. Kalp-Shaltiel, B. Lucier, I. Menache, J. S. Naor, and
J. Yaniv. Truthful online scheduling with commitments. In Proceedings
of the Sixteenth ACM Conference on Economics and Computation,
pages 715–732. ACM, 2015.
[7] P. Barham, B. Dragovic, K. Fraser, S. Hand, T. Harris, A. Ho,
R. Neugebauer, I. Pratt, and A. Warfield. Xen and the art of virtualization.
In ACM SIGOPS Operating Systems Review, volume 37,
pages 164–177. ACM, 2003.
[8] E. Boutin, J. Ekanayake, W. Lin, B. Shi, J. Zhou, Z. Qian, M. Wu, and
L. Zhou. Apollo: Scalable and coordinated scheduling for cloud-scale
computing. In OSDI, pages 285–300, Broomfield, CO, Oct. 2014.
USENIX Association.
[9] Y. Chen, A. Ganapathi, R. Griffith, and R. Katz. The case for evaluating
mapreduce performance using workload suites. In Symposium on
Modelling, Analysis, and Simulation of Computer and Telecommunication
Systems, 2011.
[10] C. Curino, D. E. Difallah, C. Douglas, S. Krishnan, R. Ramakrishnan,
and S. Rao. Reservation-based scheduling: If you’re late don’t blame
us! In SoCC, 2014.
[11] A. D. Ferguson, P. Bodik, S. Kandula, E. Boutin, and R. Fonseca.
Jockey: guaranteed job latency in data parallel clusters. In Proceedings
of the ACM European Conference on Computer Systems, EuroSys,
2012.
[12] A. Ghodsi, M. Zaharia, B. Hindman, A. Konwinski, S. Shenker, and
I. Stoica. Dominant resource fairness: Fair allocation of multiple resource
types. In NSDI, volume 11, pages 24–24, 2011.
[13] R. Grandl, G. Ananthanarayanan, S. Kandula, S. Rao, and A. Akella.
Multi-resource packing for cluster schedulers. In ACM SIGCOMM
Computer Communication Review, volume 44, pages 455–466. ACM,
2014.
[14] A. Greenberg, J. Hamilton, D. A. Maltz, and P. Patel. The cost of a
cloud: research problems in data center networks. ACM SIGCOMM
computer communication review, 2008.
[15] B. Hindman, A. Konwinski, M. Zaharia, A. Ghodsi, A. D. Joseph,
R. Katz, S. Shenker, and I. Stoica. Mesos: a platform for fine-grained
resource sharing in the data center, 2011.
[16] N. Jain, I. Menache, J. S. Naor, and J. Yaniv. A truthful mechanism
for value-based scheduling in cloud computing. Theory of Computing
Systems, 54(3):388–406, 2014.
[17] S. A. Jyothi, C. Curino, I. Menache, S. M. Narayanamurthym, A. Tumanov,
and et. al. Morpheus: Towards automated slos for enterprise
clusters. In OSDI, 2016.
[18] K. Karanasos, S. Rao, C. Curino, C. Douglas, K. Chaliparambil, G. M.
Fumarola, S. Heddaya, R. Ramakrishnan, and S. Sakalanaga. Mercury:
Hybrid centralized and distributed scheduling in large shared
clusters. In ATC, 2015.
[19] C. Kilcioglu and J. M. Rao. Competition on price and quality in cloud
computing. In WWW, 2016.
[20] I. Menache, A. Ozdaglar, and N. Shimkin. Socially optimal pricing
of cloud computing resources. In ICST Conference on Performance
Evaluation Methodologies and Tools, 2011.
[21] I. Menache, O. Shamir, and N. Jain. On-demand, spot, or both: Dynamic
resource allocation for executing batch jobs in the cloud. In
11th International Conference on Autonomic Computing (ICAC 14),
pages 177–187, 2014.
[22] P. Menage, P. Jackson, and C. Lameter. Cgroups. Available on-line at:
http://www.mjmwired.net/kernel/Documentation/
cgroups.txt, 2008.
[23] D. Merkel. Docker: lightweight linux containers for consistent development
and deployment. Linux Journal, 2014(239):2, 2014.
[24] K. Ousterhout, P. Wendell, M. Zaharia, and I. Stoica. Sparrow:
Scalable scheduling for sub-second parallel jobs. Technical Report
UCB/EECS-2013-29, EECS Department, University of California,
Berkeley, Apr 2013.
[25] J. Rasley, K. Karanasos, S. Kandula, R. Fonseca, M. Vojnovic, and
S. Rao. Efficient queue management for cluster scheduling. In Proceedings
of the Eleventh European Conference on Computer Systems,
page 36. ACM, 2016.
[26] D. Sarkar. Introducing hdinsight. In Pro Microsoft HDInsight, pages
1–12. Springer, 2014.
[27] B. Sharma, R. K. Thulasiram, P. Thulasiraman, S. K. Garg, and
R. Buyya. Pricing cloud compute commodities: a novel financial economic
model. In IEEE-CCGRID, 2012.
[28] Y. Song, M. Zafer, and K.-W. Lee. Optimal bidding in spot instance
market. In INFOCOM, 2012 Proceedings IEEE, pages 190–198.
IEEE, 2012.
[29] M. A. team. Azure batch: Cloud-scale job scheduling and compute
management. In https://azure.microsoft.com/en-us/
services/batch/, 2015.
[30] A. Tumanov, T. Zhu, J. W. Park, M. A. Kozuch, M. Harchol-Balter,
and G. R. Ganger. Tetrisched: global rescheduling with adaptive planahead
in dynamic heterogeneous clusters. In Eurosys, 2016.
[31] V. K. Vavilapalli, A. C. Murthy, C. Douglas, S. Agarwal, M. Konar,
R. Evans, T. Graves, J. Lowe, H. Shah, S. Seth, et al. Apache hadoop
yarn: Yet another resource negotiator. In ACM - SoCC, 2013.
[32] A. Velte and T. Velte. Microsoft virtualization with Hyper-V.
McGraw-Hill, Inc., 2009.
[33] A. Verma, L. Pedrosa, M. Korupolu, D. Oppenheimer, E. Tune, and
J. Wilkes. Large-scale cluster management at google with borg. In
Eurosys, 2015.
[34] C. A. Waldspurger. Memory resource management in vmware esx
server. SOSP, 2002.
[35] M. Zaharia, D. Borthakur, J. Sen Sarma, K. Elmeleegy, S. Shenker,
and I. Stoica. Delay Scheduling: A Simple Technique for Achieving
Locality and Fairness in Cl
