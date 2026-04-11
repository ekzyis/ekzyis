---
title: Fixing A High-Impact Bug In Saving Satoshi
date: 2025-12-03
tags: web dev
---

Around the same time [I resigned from Stacker
News](https://stacker.news/items/1309092?commentId=1310102),
[Jonas](https://adamjonas.com/) from [Chaincode Labs](https://chaincode.com/)
posted about the [BOSS Challenge](https://bosschallenge.xyz/) on [Stacker
News](https://stacker.news/items/1299187). That was very fortunate timing, since
I was indeed looking for opportunities to continue working in bitcoin,
especially as a grantee.

To join the program starting on January 12, 2026, one had to complete [Saving
Satoshi](https://savingsatoshi.com/). It wasn't clear whether everyone who
finished Saving Satoshi would get invited, so I rushed to complete it ASAP. I
hoped it would indicate "this guy REALLY wants to get in", since it was true. I
heard a lot of good stuff about Chaincode Labs, including their former annual
residencies, which this online program replaced.

After the second challenge, I noticed [**this
bug**](https://github.com/saving-satoshi/saving-satoshi/issues/1319) in the
desktop layout of the website that stopped me from proceeding. I found a way
around it by switching to the mobile layout. However, I knew that this was a
serious bug. If it stops me from proceeding, it will also stop others, and they
might not find a way around it like I did.

I started to look into it. Especially as a former web developer myself, I could
empathize with the team. These bugs are easy to miss, but have a huge impact. You
would need to use your product _all the time_ to never miss them. If someone
showed me such a bug in my code, it would be an emergency.

I saw it as a great opportunity to show my skills relevant for BOSS: I get stuff
done on my own, I demonstrate that I understood a problem well, including its
context and solutions, and I have empathy. And I would say I succeeded!
[Satsie](https://github.com/satsie), a maintainer of the website, acknowledged
it:

> THANK YOU [@ekzyis](https://github.com/ekzyis)! That bug you found was a
> showstopper and needed to be fixed immediately. Really appreciate this
> detailed bug report and all the videos you included. You made it very easy to
> understand what was going on and verify your fix. I also like that you added
> the two other solutions you investigated, with corresponding video.
>
> Congratulations on your first contribution to Saving Satoshi! 🎉

Thank you Satsie!
