import os
from django.contrib.auth import get_user_model
from django.utils import timezone
from datetime import timedelta
from oauth2_provider.models import AccessToken
from cookbook.models import Recipe, Space, UserSpace, UserPreference
from django_scopes import scope, scopes_disabled
from django.contrib.auth.models import Group

User = get_user_model()
with scopes_disabled():
    Recipe.objects.all().delete()
    UserSpace.objects.all().delete()
    Space.objects.all().delete()
    User.objects.all().delete()

user, _ = User.objects.get_or_create(username='admin', email='admin@example.com')
user.set_password('admin')
user.is_superuser = True
user.is_staff = True
user.save()

from oauth2_provider.models import AccessToken, Application

# Create an Application for the token
app, _ = Application.objects.get_or_create(
    name='Tandoor MCP E2E',
    user=user,
    client_type=Application.CLIENT_CONFIDENTIAL,
    authorization_grant_type=Application.GRANT_PASSWORD,
)

access_token, created = AccessToken.objects.get_or_create(
    user=user,
    token="e2e_test_token",
    defaults={
        "application": app,
        "expires": timezone.now() + timedelta(days=3650),
        "scope": "read write"
    }
)

from cookbook.helper.permission_helper import create_space_for_user

# Tandoor requires a Space
us = create_space_for_user(user, 'e2e_test_space')
us.active = True
us.save()
space = us.space

# Create a test recipe in the space
with scope(space=space):
    Recipe.objects.get_or_create(name='Tandoori Chicken', created_by=user, space=space, defaults={'description': 'E2E test'})

print("Successfully seeded E2E user, token, space, and recipe.")
