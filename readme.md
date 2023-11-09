# 说明

这是什么？

xmail 是一个非常简易的程序，来帮助我们搭建一个接受邮件的服务，他借助 cloudflare email worker 接受邮件，使用 cloudflare d1 缓存邮件， 然后通过 xmail 来拉取有 cloudflare worker 提供的 api 获取邮件并缓存到本地 mongo 上进行查看。由 xmail 提供的服务进行查看。

因此你可以搭建属于自己的邮件接收器，将来自不同地方的邮件进行汇集。同事可以通过一些监听，完成更好的脚本处理能力。
