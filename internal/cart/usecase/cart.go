package usecase

import (
	"context"
	"sort"
	"sync"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"

	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUseCase) Update(ctx context.Context, sc models.Scope, input cart.UpdateInput) (cart.UpdateOutput, error) {
	err := uc.validateCartItem(ctx, sc, input.NewItemList)
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Update.validateCartItem", err)
		return cart.UpdateOutput{}, err
	}

	var data getDataOutput
	if len(input.NewItemList) > 0 {
		data, err = uc.getDataCartItems(ctx, sc, input.NewItemList)
		if err != nil {
			uc.l.Errorf(ctx, "cart.Usecase.Update.getDataCartItems", err)
			return cart.UpdateOutput{}, err
		}
	}

	var shopIDsSet = util.RemoveDuplicates(data.ShopIDs)
	carts, err := uc.repo.ListCart(ctx, sc, cart.ListOption{
		CartFilter: cart.CartFilter{
			ShopIDs: shopIDsSet,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Update.ListCart", err)
		return cart.UpdateOutput{}, err
	}

	if len(input.NewItemList) > 0 && len(carts) != len(shopIDsSet) {
		uc.l.Errorf(ctx, "cart.Usecase.Update.ListCart", cart.ErrCartNotFound)
		return cart.UpdateOutput{}, cart.ErrCartNotFound
	}
	var wg sync.WaitGroup
	var wgErr error
	var Mutex sync.Mutex
	var updatedCarts []models.Cart

	for _, c := range carts {
		wg.Add(1)
		go func(c models.Cart) {
			defer wg.Done()
			for i, item := range data.CartItems {
				if item.Quantity == 0 {
					data.CartItems = append(data.CartItems[:i], data.CartItems[i+1:]...)
				}
			}

			if len(data.CartItems) == 0 {
				err := uc.repo.Delete(ctx, sc, c.ID)
				if err != nil {
					uc.l.Errorf(ctx, "cart.Usecase.Update.Delete", err)
					wgErr = err
					return
				}
				return
			}

			updatedCart, err := uc.repo.Update(ctx, sc, cart.UpdateCartOption{
				Model:       c,
				NewItemList: data.CartItems,
			})
			if err != nil {
				uc.l.Errorf(ctx, "cart.Usecase.Update.Update", err)
				wgErr = err
				return
			}
			Mutex.Lock()
			updatedCarts = append(updatedCarts, updatedCart)
			Mutex.Unlock()
		}(c)
	}

	wg.Wait()
	if wgErr != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Update.Update", wgErr)
		return cart.UpdateOutput{}, wgErr
	}

	return cart.UpdateOutput{
		Carts: updatedCarts,
		Shops: data.Shops,
	}, nil
}

func (uc implUseCase) Add(ctx context.Context, sc models.Scope, input cart.CreateCartInput) error {
	err := uc.validateCartItem(ctx, sc, []cart.CartItemInput{
		{
			ProductID: input.ProductID,
			Quantity:  input.Quantity,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Add.validateCartItem: %v", err)
		return err
	}

	data, err := uc.getDataCartItems(ctx, sc, []cart.CartItemInput{
		{
			ProductID: input.ProductID,
			Quantity:  input.Quantity,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Add.getDataCartItems: %v", err)
		return err
	}

	existingCart, err := uc.repo.GetOne(ctx, sc, cart.GetOneOption{
		CartFilter: cart.CartFilter{
			ShopIDs: []string{data.ShopIDs[0]},
		},
	})
	if err == mongo.ErrNoDocuments {
		_, err = uc.repo.Create(ctx, sc, cart.CreateCartOption{
			ShopID: data.ShopIDs[0],
			CartItemList: []models.CartItem{
				{
					ProductID: data.CartItems[0].ProductID,
					Quantity:  input.Quantity,
				},
			},
		})
		if err != nil {
			uc.l.Errorf(ctx, "cart.Usecase.Add.Create: %v", err)
			return err
		}

		return nil
	} else if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Add.GetOne: %v", err)
		return err
	}

	var found bool
	var newItems []models.CartItem
	existCartProductIDs := []string{}

	for _, item := range existingCart.Items {
		existCartProductIDs = append(existCartProductIDs, item.ProductID.Hex())
	}

	products, err := uc.shopUc.ListProduct(ctx, models.Scope{}, shop.ListProductInput{
		GetProductFilter: shop.GetProductFilter{
			IDs: existCartProductIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Add.ListProduct: %v", err)
		return err
	}

	invensMap := make(map[string]models.Inventory)
	for _, p := range products.Products {
		invensMap[p.P.ID.Hex()] = p.Inventory
	}

	for _, item := range existingCart.Items {
		if item.ProductID == data.CartItems[0].ProductID {
			item.Quantity += input.Quantity
			found = true
		}

		if err := uc.checkStock(ctx, sc, invensMap[item.ProductID.Hex()], item.Quantity); err != nil {
			uc.l.Errorf(ctx, "cart.Usecase.Add.checkStock", err)
			return err
		}

		newItems = append(newItems, item)
	}

	if !found {
		newItems = append(newItems, models.CartItem{
			ProductID: data.CartItems[0].ProductID,
			Quantity:  input.Quantity,
		})
	}

	_, err = uc.repo.Update(ctx, sc, cart.UpdateCartOption{
		Model:       existingCart,
		NewItemList: newItems,
	})
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Add.Update", err)
		return err
	}

	return nil
}

func (uc implUseCase) GetCart(ctx context.Context, sc models.Scope, opt cart.GetOption) (cart.GetCartOutput, error) {
	carts, pag, err := uc.repo.GetCart(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "cart.usecase.GetCart: %v", err)
		return cart.GetCartOutput{}, err
	}

	var productItem cart.ProductItem
	var getCartItems []cart.GetCartItem
	var cartProductQuantityMap = make(map[string]map[string]int)
	var cartProductMap = make(map[string][]string)
	var cartShopMap = make(map[string]models.Shop)
	var shopIDs []string

	for _, v := range carts {
		shopIDs = append(shopIDs, v.ShopID.Hex())
		cartProductQuantityMap[v.ID.Hex()] = make(map[string]int)
		for _, item := range v.Items {
			cartProductMap[v.ID.Hex()] = append(cartProductMap[v.ID.Hex()], item.ProductID.Hex())
			cartProductQuantityMap[v.ID.Hex()][item.ProductID.Hex()] = item.Quantity
		}
	}

	listShops, err := uc.shopUc.ListShop(ctx, models.Scope{}, shop.GetShopsFilter{
		IDs: shopIDs,
	})
	if err != nil {
		uc.l.Errorf(ctx, "cart.usecase.GetCart: %v", err)
		return cart.GetCartOutput{}, err
	}

	for _, v := range listShops {
		cartShopMap[v.ID.Hex()] = v
	}

	for _, v := range carts {
		var productItems []cart.ProductItem
		var carProductMediaMap = make(map[string][]models.Media)

		listProducts, err := uc.shopUc.ListProduct(ctx, models.Scope{}, shop.ListProductInput{
			GetProductFilter: shop.GetProductFilter{
				IDs: cartProductMap[v.ID.Hex()],
			},
		})
		if err != nil {
			uc.l.Errorf(ctx, "cart.usecase.GetCart: %v", err)
			return cart.GetCartOutput{}, err
		}

		for _, p := range listProducts.Products {
			productItem = cart.ProductItem{
				ProductID:   p.P.ID.Hex(),
				Medias:      p.Images,
				Quantity:    cartProductQuantityMap[v.ID.Hex()][p.P.ID.Hex()],
				ProductName: p.P.Name,
				Price:       p.P.Price,
			}
			carProductMediaMap[p.P.ID.Hex()] = p.Images
			productItems = append(productItems, productItem)
		}

		if len(productItems) > 0 {
			getCartItems = append(getCartItems, cart.GetCartItem{
				Cart:                v,
				CartProductMediaMap: carProductMediaMap,
				Products:            productItems,
				Shop:                cartShopMap[v.ShopID.Hex()],
			})
		}
	}

	sort.Slice(getCartItems, func(i, j int) bool {
		return getCartItems[i].Cart.CreatedAt.After(getCartItems[j].Cart.CreatedAt)
	})

	return cart.GetCartOutput{
		CartOutPut: getCartItems,
		Paginator:  pag,
	}, nil
}
