package server

import (
	"context"
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	useroperatorv1alpha1 "github.com/lcpu-club/user-operator/api/v1alpha1"
)

func (s *Server) initKube() (err error) {
	s.kubeClient, err = dynamic.NewForConfig(&rest.Config{
		Host: s.upstream.String(),
		TLSClientConfig: rest.TLSClientConfig{
			CAFile: s.sm.GetCAFile(),
		},
		BearerToken:     s.sm.GetToken(),
		BearerTokenFile: s.sm.GetTokenFile(),
		QPS:             20.0,
		Burst:           30,
	})
	if err != nil {
		return err
	}

	s.scheme = runtime.NewScheme()
	err = useroperatorv1alpha1.AddToScheme(s.scheme)
	if err != nil {
		return err
	}

	s.userGVR = useroperatorv1alpha1.GroupVersion.WithResource("users")
	s.userGVK = useroperatorv1alpha1.GroupVersion.WithKind("User")

	return nil
}

func (s *Server) kubeGetUser(uid string) (*useroperatorv1alpha1.User, error) {
	obj, err := s.kubeClient.Resource(s.userGVR).Get(context.TODO(), uid, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	u := &useroperatorv1alpha1.User{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Server) kubeCreateUser(u *useroperatorv1alpha1.User) (*useroperatorv1alpha1.User, error) {
	u.SetGroupVersionKind(s.userGVK)
	unstruct, err := runtime.DefaultUnstructuredConverter.ToUnstructured(u)
	if err != nil {
		return nil, err
	}

	obj, err := s.kubeClient.Resource(s.userGVR).Create(context.TODO(), &unstructured.Unstructured{
		Object: unstruct,
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("Hello 3 %v", err)
	}

	ur := &useroperatorv1alpha1.User{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, ur)
	if err != nil {
		return nil, err
	}

	return ur, nil
}

func (s *Server) reconcileUser(ii *ImpersonateInfo) error {
	u := &useroperatorv1alpha1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: ii.UID,
		},
		Spec: useroperatorv1alpha1.UserSpec{
			Username: ii.Username,
			UID:      ii.UID,
			Groups:   ii.Group,
			Extra:    ii.Extra,
		},
	}

	u.SetGroupVersionKind(s.userGVK)

	_, err := s.kubeGetUser(ii.UID)

	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			_, err = s.kubeCreateUser(u)
			log.Println("Created user", ii.UID)
			return err
		}
		return err
	}

	return nil
}
