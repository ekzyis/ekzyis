---
title: Lightning Prediction Market MVP - delphi.market
date: 2023-12-04
tags: bitcoin
sn_id: 337637
---

This is the announcement that I announced
[here](https://stacker.news/items/336740):

---

## **I finally got the MVP done for a [prediction market on lightning](https://delphi.market/).**


_If you don't know what a prediction market is, watch [this
video](https://www.youtube.com/watch?v=xA27x7GRMZQ)._

---

## But first and foremost: **To prevent loss of funds, the market is running on [mutinynet](https://blog.mutinywallet.com/mutinynet/) and not mainnet until I have enough confidence in the code for mainnet. This means the sats have no value and you need to use the [mutiny faucet](https://faucet.mutinywallet.com/) to pay invoices and the [mutiny wallet (signet)](https://signet-app.mutinywallet.com/) to test withdrawals as mentioned in [/about](https://delphi.market/#/about).**

**Thanks to [@benthecarman](https://stacker.news/benthecarman)
[@TonyGiorgio](https://stacker.news/TonyGiorgio) and Paul (I don't know his nym
on here, lol) for mutinynet and the provided tools! Less stuff I had to worry
about on my own.**

---

## **Features**

### 1. **create your own binary markets**

On the main page, you can immediately see a button to create a new market.
You'll need to login with lightning though. You can use any LNURL-auth capable
wallet for that. LNURL-auth is independent from which network is used for
transactions :)

For now, you can only create binary markets. This means that you can only bet on
events which outcome is either `YES` or `NO`:

![](./dm001.webp)
market creation form

Be aware that you can currently not change the description or delete the market
after creation. In the future, I might allow edits for a few minutes after
creation.

### 2. **trading of YES and NO shares**

Basically the heart of the market. As a market would be useless without trading,
you can obviously create BUY and SELL orders for YES and NO shares. **Every
order must be matched with an order of a different market participant.**

For example, if you want to buy 2 YES shares at a price of 70 sats per share or
short: `BUY 2 YES @ 70 sats`, your order must match a `BUY 2 NO @ 30 sats`
order.

![](./dm002.webp)
buy order form

The interface might still be overwhelming even though I iterated on it several
times already.

The input labeled with `how much?` specifies how many sats you want to bet. With
`how sure?` you specify the price per share since the share prices in a
prediction market reflect what the market currently thinks how likely an outcome
is (at least in theory).

### 3. **interactive order book**

To match existing orders easily, you can simply click on `PENDING` and a context
menu will show - assuming it's not your own order, else it will show a button to
cancel your order:

![](./dm003.webp)
interactive order book

On click on `match`, you'll get redirected back to the order form but with
prefilled values. If you submit this order, this will match the other order
after payment. For example, if you click on `match` in the screenshot, you would
get redirected to this order form with `BUY 2 YES @ 20 sats` prefilled:

![](./dm004.webp)
after clicking on match

### 4. **basic market stats**

On the `stats` page, you can see the total sats locked in the market (volume)
and how many are on which side. `YES` 75% means that 75% of all sats were put
into YES share which means that the market thinks that the chance that YES will
win is 75%. However, tbh, I am not entirely sure if the chart works, lol

![](./dm005.webp)
stats page

### 5. **withdrawal**

You can withdraw any sats you received via payout or canceling an order (you'll
get a refund for any invoice you paid for the order) in your wallet:

![](./dm006.webp)
wallet page

![](./dm007.webp)
withdrawal via bolt11 payment request

### 6. **list of all invoices**

Not much to say here. Just a list of all your invoices with their status so you
can pay them later or keep track of what you've done:

![](./dm008.webp)
list of all your invoices

### 7. **list of all orders**

Similar to the invoice list. You can also cancel your orders from this list (and
not only from within a market order book) by clicking on `PENDING` to show the
context menu.

![](./dm009.webp)
list of all your orders

### 8. **market settlement**

If you created a market, you can decide the outcome of the market via a hidden
`settings` page:

![](./dm010.webp)
market settings page

This is obviously not ideal and unlike what I initially mentioned: that the
market creator (or me, as the site operator) would only be used as a fallback if
the trading parties don't find consensus on the outcome. However, this idea
didn't make it into the MVP.

If anyone even remembers, I also mentioned that I would be using HODL invoices
to wait with settling the invoices until a counter party was found. If no
counter party was found, their sats would get returned for free - no network
fees. That also didn't make it into the MVP and I am no longer sure that's even
something users would want since it trades UX for a potentially insignificant
risk of having to pay the lightning network again for a withdrawal if you found
no counter party.

---

## **Planned**

- actually show end date that was picked during market creation to user
- PWA with push notifications
- fix chart (?)
- code cleanup / refactor (that's probably always a thing, lol)
- tests, tests, tests
- other stuff that I forgot now
- ...
- mainnet

---

That's it for now. I want to go to sleep and writing all of this took some time
that I was worried [about my account deletion
deadline](https://stacker.news/items/336740), lol

Now go, find all the bugs that I have hidden for you! Feel free to create new
markets. I have created some orders for the existing markets so you can use them
to immediately find some matches :)

Hopefully I'll wake up to some interesting prediction markets - or funny :)

---

_some somewhat random pings: [@Natalia](https://stacker.news/Natalia)
[@grayruby](https://stacker.news/grayruby)
[@Undisciplined](https://stacker.news/Undisciplined)
[@nemo](https://stacker.news/nemo) [@siggy47](https://stacker.news/siggy47)
[@orthwyrm](https://stacker.news/orthwyrm)
[@elvismercury](https://stacker.news/elvismercury)
[@carlosfandango](https://stacker.news/carlosfandango)_