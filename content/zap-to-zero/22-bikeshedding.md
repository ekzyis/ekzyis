---
title: Zap to Zero Day 22 | Bikeshedding
date: 2024-01-19
sn_id: 394081
---

![Leonardo_Diffusion_XL_A_sleek_futuristic_bicycle_gleaming_in_t_0.jpg](https://m.stacker.news/13033)

_A sleek, futuristic bicycle gleaming in the harsh light of a post-apocalyptic world, its expensive
price tag a stark contrast to the decaying nuclear power plant looming in the background._

> _Bikeshedding, also known as Parkinson’s law of triviality, describes our tendency to devote a
> disproportionate amount of our time to menial and trivial matters while leaving important matters
> unattended._

— The Decision Lab, [_What is Bikeshedding?_](https://thedecisionlab.com/biases/bikeshedding)

---

Yesterday, I worked on integrating [Nostr Wallet
Connect](https://github.com/nostr-protocol/nips/blob/master/47.md) into SN. While I was trying to
figure out how I can check if the connection URI from the wallet is working—don't want to throw
preventable errors in your faces on zap attempts[^1]—I found this comment from @benthecarman that
made me laugh hard (still laughing every time I read it):

>> Maybe this is a new nip indeed. Requesting invoices shouldn't require authorization. The need
>> looks more like a data vending machine that replies when a key is pinged than a nwc action.
>
> you guys bike shedding so much more than needed. the wallet connect functionality should have the
> ability to receive, that is all this does

— @benthecarman, [_NIP-47: Nostr Wallet Connect
Extensions_](https://github.com/nostr-protocol/nips/pull/685#issuecomment-1654968901), Jul 28, 2023

And the discussion only got more funny:

>> As it stands, the change alone doesn't seem to be a replacement for the public lnurlp resolution
>> when asking for an invoice from somebody else's wallet.
>
> you can build zaps on top of this, you can build a lightning address server on top of this, you
> can build a webstore on top of this.
>
> the point of this is to allow the functionality of requesting invoices, not building a specific
> app/use case like zaps

However, people weren't satisfied that this PR doesn't solve world hunger—yet. So more and more
people asked if it would make sense to add more and more things [^2].

At some point, @supertestnet [joined
in](https://github.com/nostr-protocol/nips/pull/685#issuecomment-1668076555) and asked about the
status of this PR since it wasn't clear.

>>> It looks like this spec-change has been approved and is now live in a popular NWC
>>> implementation. Can it be merged so that other client devs know how to interoperate? This would
>>> help ensure different implementations use the same api for the same functionality, and it would
>>> help client devs like me know what api methods to call if they want to use this functionality
>>
>> need 2 implementations
>
> ok
>
> looks like a second implementation is on its way, hopefully it is here soon
>
> [getAlby/nostr-wallet-connect#112
> (comment)](https://github.com/getAlby/nostr-wallet-connect/issues/112#issuecomment-1668903161)

After this was done, @benthecarman asked on Sep 3, 2023:

> https://github.com/getAlby/nostr-wallet-connect now implements this, good to merge?

But we weren't done with @benthecarman yet. Even @TonyGiorgio [joined the
fun](https://github.com/nostr-protocol/nips/pull/685#issuecomment-1728100073):

> How about a keysend method? That would enable wallets to be the funding source for a lot of value
> for value podcast stuff.

After @benthecarman also added this method, he dared to ask if it was good to merge now:

> I think this PR is getting pretty big at this point. We should have enough functionality now for
> any kind of wallet.
>
> Good to merge? [@bumi](https://github.com/bumi) [@rolznz](https://github.com/rolznz)

But since new methods were added to the spec, these new methods also had to be implemented now:

> @benthecarman we haven't implemented the new methods yet, if possible it would be great to do this
> first in case we notice any changes that need to be made.

On Oct 11, 2023, `get_info` [was proposed as the final
method](https://github.com/nostr-protocol/nips/pull/685#issuecomment-1756752264). This was the
method I was currently trying to use to check the connection status.

>  I would like to propose one final additional method - what do you guys think?
>
> **Get Info**
>
> Returns information about the user's node

The discussion continued. @Semisol also[ joined in and did a
review](https://github.com/nostr-protocol/nips/pull/685#discussion_r1376599737) on Oct 30, 2023. But
since the PR has deviated a lot since it started as a simple way to add receiving to NWC, they were
confused:

>> I don't get the goal here.
>
> The original intention was to be able to create invoices over NWC. That way NWC can be used as a
> replacement for things like LNURL in the future.
>
> The alby guys wanted to add full functionality to be able to get all other payment information so
> we added list payments and invoices.

The latest comments were about `multi_pay_keysend` and `multi_pay_invoice` which I think is related
to [_Atomic Multi-path Payments
(AMP)_](https://docs.lightning.engineering/lightning-network-tools/lnd/amp).

But I recommend that you checkout the discussion yourself and see how the ~sausage~ bleeding edge of
technology is made!

---

## [Satistics](https://stacker.news/satistics?inc=stacked,spent)

| Date | Spent | Stacked (Rewards) | Posts | Comments | Rewarded |
| --- | --- | --- | --- | --- | --- |
| [2023-12-28](https://stacker.news/items/370188) | 13k | 8808 (n/a) | 2 | 35 | n/a  |
| [2023-12-29](https://stacker.news/items/371370) | 16.1k | 15.6k (5222) | 3 | 52 | ⚡ |
| [2023-12-30](https://stacker.news/items/372534) | 10.8k | 9752 (7026) | 1 | 41 | ✍️ |
| [2023-12-31](https://stacker.news/items/373620) | 20.5k | 17.9k (4379) | 5 | 61 | ⚡ |
| **[2024-01-01](https://stacker.news/items/374881)** | 12.5k | 10.7k (7684) | 3 | 47 | ✍️ |
| [2024-01-02](https://stacker.news/items/376287) | 16k | 19.5k (9353) | 6 | 36 | ✍️ |
| [2024-01-03](https://stacker.news/items/377500) | 15.9k | 15.6k (6729) | 2 | 46 | ⚡ |
| [2024-01-04](https://stacker.news/items/378870) |	11.4k | 11.4k (~3954~ ~4093~ 4131) | 3 | 38 | ✍️ |
| [2024-01-05](https://stacker.news/items/379663) |	11.3k | 11.4k (~3954~ 4092) | 1 | 41 | ? |
| [2024-01-06](https://stacker.news/items/379663) | 6691 | 6282 (~3665~ 3954) | 0 | 38 | ✍️ |
| [2024-01-07](https://stacker.news/items/380810) | 8053 | 8503 (~1219~ 3665) | 3 | 20 | ✍️ |
| **[2024-01-08](https://stacker.news/items/382187)** | 8873 | 9164 (1219) | 2 | 12 | ⚡ |
| [2024-01-09](https://stacker.news/items/383365) | 5828 |  6808 (4649) | ~2~ 6 | ~34~ 35 | ✍️ |
| [2024-01-10](https://stacker.news/items/384428) | 14.1k | 14.4k (4857) | 3 | 22 | ⚡ |
| [2024-01-11](https://stacker.news/items/385785) | 11.8k | 10.4k (4109) | 3 | 22 | ✍️ |
| [2024-01-12](https://stacker.news/items/386849) | 8743 | 8016 (4778) | 3 | 41 | ✍️ |
| [2024-01-13](https://stacker.news/items/387537) | 9393 | 9339 (3116) | 2 | 17 | ⚡ | 
| [2024-01-14](https://stacker.news/items/388997) | 14.2k | 6697 (3533) | 4 | 41 |  ✍️ |
| [2024-01-15](https://stacker.news/items/390374) | 10.2k | 11.3k (3395) | 1 | 15 | ⚡ |
| **[2024-01-16](https://stacker.news/items/391205)** | 6867 | 11.1k (2500) | 2 | 27 | ⚡ |
| [2024-01-17](https://stacker.news/items/392798) | 6008 | 6931 (3982) | 1 | 21 | ⚡ |
| [2024-01-18]() | 6827 | 7606 (2544) | 2 | 27 | ⚡ |
| 2024-01-19 | TBD | TBD (3755) | TBD | TBD | ⚡|

![2024-01-19-191147_605x389_scrot.png](https://m.stacker.news/13049)

![2024-01-19-191520_906x235_scrot.png](https://m.stacker.news/13052)

![Figure_2.png](https://m.stacker.news/13053)

![Figure_3.png](https://m.stacker.news/13054)

I turned my PC off yesterday. I usually never turn it off so I can easily continue with whatever I
was doing. But yesterday, I felt like I need to turn it off to help me go to sleep.

Since the script that polled my SN balance every 5 second was running on it, it died. I realized my
mistake just a few minutes ago and restarted it. Need to move it to a server of mine. This should be
one of those 5 minutes things that you should just immediately do instead of thinking about _when_
to do. But I really want to finish this post first. And then I will have probably forgotten about
it.

Btw, I am actually quite surprised that it took so long for something like this to happen. This is
the code that was running without any hiccups since [Day 12](https://stacker.news/items/382187):

```
#!/usr/bin/env bash

COOKIE="<INSERT COOKIE HERE>"

ENDPOINT=https://stacker.news/api/graphql

while true;
do
  sats=$(curl \
    -X POST \
    -H "Content-type: application/json" \
    -H "Cookie: $COOKIE" \
    --data '{ "query": "query ME { me { privates { sats } id } }" }' \
    --silent "$ENDPOINT" \
  | jq '.data.me.privates.sats')
  line="$(date +%s) $sats"
  echo "$line" | tee -a me_sats
  sleep 5
done
```

---

## Recent Superzaps

### 1. [Confessions of a burnt-out father](https://stacker.news/items/304937)

A while ago, @cryotosensei shared his pain of being a father of two kids with us. Pain that he
expected but pain that still hurts.

I don't want to quote too much since I don't want people to feel like they don't have to read the
actual post since they read "my summary". But this isn't supposed to be a summary. This is supposed
to make you interested in reading the actual post and zap it. Consider this to be the abstract of a
scientific paper. If you only read the abstract of the [Bitcoin
Whitepaper](https://stacker.news/bitcoin.pdf), did you really read the Bitcoin Whitepaper? The
answer should be obvious.

So here is my favorite line where I felt a real connection to @cryotosensei even though I am far
away from being a parent—and I am not sure if I ever _really_ want to or will have kids.

> I don’t have to be around people all the time. I’m perfectly fine signing up for seminars and
> workshops and attending them alone. I miss being exposed to people outside my usual realm of work
> and expanding my horizons. I hate that I have to rush home every week day after work because I
> need to pick up my kids from the childcare centre. I hate that I have to spend my weekends taking
> my boy to the indoor playground so that he can release his pent-up energies. In fact, right now
> I’m at the indoor playground. All I see around me are equally bored parents passing the time
> buried in their handphones. Since when did my activity-filled life turn out to be so mundane?

So please, check out the post from @cryotosensei and share whatever you have to share with him. I
think he deserves it. Even though he sometimes talk a little too much for my liking about his wife
being Japanese, lol.

---

### 2. [What’s your SN routine like?](https://stacker.news/items/393176)

Another great post from @cryotosensei. This one is from yesterday.

I loved how @cryotosensei not just asked us about our routine ([like a
fed](https://stacker.news/items/392090)) but actually shared his own routine first (like a very
sophisticated fed). That's what SN is about: Value 4 Value!

> I have had the cowboy hat in my possession for 24 days, so suffice it to say that posting and
> commenting here on Stacker News have become a comfortable routine.
>
> While other regulars fiddle around with the optimal amount of zaps they should give, I have been
> more preoccupied with incorporating the use of SN into my life. So after 24 days, this is
> basically my routine these days:
>
> [...]
>
> What’s your SN routine like?

I also loved reading the comments. @cryotosensei basically created a good prompt to create bonds
between old cowboys and people learning the ropes. For example, @siggy47, the founder of
~bitcoin_beginners, replied this:

> Can I just say that I see you all the time making quality contributions here. I think you may be
> cheating yourself of sats with your method. You will earn far more by tipping more. Build up some
> more sats, then tip quality posts early. For more details, check @Natalia's recent post. It really
> pays to pay others before yourself, as crazy as that sounds.
As an example, the first thing I do in the morning is donate a load of sats to the site. It's not
completely out of the goodness of my heart. First of all, then I don't need to worry about my cowboy
hat. Secondly, I know I have a chance to earn it back with good tipping and posting throughout the
day.

— @siggy47, [#393198](https://stacker.news/items/393198)

_Maybe there should also be ~meta_beginners territory? Half-joking._

---

### 3. [The Mysterious Life And Death Of Mircea Popescu](https://stacker.news/items/204836)

In SN terms, this is a very old post. When @siggy47 found out that some dude named Popescu died, he
dug deeper (not into his grave fortunately—at least I hope so).

> I once took a vacation with my family to Costa Rica’s pacific coast, not far from the place where
> Mircea Popescu reportedly died. It’s a beautiful part of the world, and it’s on my short list of
> places to relocate to if I get the motivation to leave the United States. I say it’s where he
> reportedly died, because there is little about Mircea Popescu that is certain. That’s the way he
> wanted it. Yes, I used the past tense, and I will keep doing so for the sake of consistency, but
> I’m not sure he’s dead.
>
> I must confess that when the news of the drowning began appearing in bitcoin publications in late
> June of 2021 I had never heard of him. He was apparently an OG bitcoiner, a major whale and a
> famous blogger.  The headlines said that he may have died holding over one million bitcoin, making
> him the only “bitcoin millionaire” other than Satoshi Nakamoto. I read elsewhere that his stash
> was not nearly that large.  It was all based on speculation. So far as I can tell, none of his
> wallets have ever been traced or identified. Everyone just took his word for it when he boasted of
> his enormous hoard. I discovered that he was controversial. He was considered a brilliant bitcoin
> philosopher and a savvy entrepreneur, while at the same time he was described as hateful, racist,
> misogynistic, and anti-semitic.  I am a bit of a history junkie, so I wanted to learn more about
> him.

I must confess that I didn't really look much into this Popescu dude until @siggy47 mentioned
yesterday that my writing style is similar; just more friendly.

> I don't want to dump more work on you, but when you're done with this experiment I hope you keep
> writing these posts everyday as a sort of blog. Your writing style is addicting. Kind of a more
> friendly Popescu. I read it to see what's going on in that head of yours more than for this
> experiment. I bet I'm not the only one.

— @siggy47, [#392821](https://stacker.news/items/392821)

I looked into [trilema.com](http://trilema.com/) and even the domain name was something that
inspired me. The front page also loads random posts on every site fresh. Also something that I can
imagine myself doing.

However, I would never host a site without TLS enabled. I know that it might not make _total_ sense
for something like a blog[^3] but with my background in cryptography, that's just a no-go for me
when it comes to my own sites.

I am very happy that @siggy47 mentioned this. I am still processing what it means for me though.

It feels like I passed some kind of writing event horizon.

---

## Challenge of the Day

I don't know, why are you asking me? Go figure out your own challenge of the day.

---

## Song of the Day

![](https://www.youtube.com/watch?v=sZmK2ZWBxQc)

> And I feel tired \
> Empty and hollow \
> Heart broken inside \
> And I feel this life \
> Has nothing for me anymore
>
> And I feel revived \
> Sacred and honored \
> One of a kind \
> And I feel this life \
> Is something I was chosen for

_Crazy to think about that I tried to learn this song on guitar 10 years ago. You can't believe how
proud I was that I figured out how to play the first 9 seconds. Not in the same speed, but the intro
was recognizable._

---

Maybe I should have ended this series while I had the chance to end on a nice number like 21. I also
created a [freebie comment](https://stacker.news/items/392840) yesterday on accident. This was one
of the original loss conditions for this series that [started as a
game](https://stacker.news/items/369003).

---

_images generated using [leonardo.ai](https://leonardo.ai/)_

[^1]: Jim Shore, [_Fail Fast_](https://www.martinfowler.com/ieeeSoftware/failFast.pdf)

[^2]: To be fair, they were all _somewhat_ related to receiving afaict. For example, `list_invoices`
    (to list invoices) and `list_payments` (to list payments) were mentioned. Later, another method
    was discovered: `list_transactions` should do what `list_invoices` and `list_payments` does in
    one call.

[^3]: Even Paul Graham's blog [didn't use TLS in the
    past](https://news.ycombinator.com/item?id=27196917). [It now does](https://paulgraham.com/) for
    some reason though.

